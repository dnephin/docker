package main

import (
	"fmt"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// TestRemoveFailsWhenImageIsInUse (with Force and without)
// check image tag still exists
func (s *DockerSuite) TestRmiWithContainerFails(c *check.C) {
	errSubstr := "is using it"

	// create a container
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")

	cleanedContainerID := strings.TrimSpace(out)

	// try to delete the image
	out, _, err := dockerCmdWithError("rmi", "busybox")
	// Container is using image, should not be able to rmi
	c.Assert(err, checker.NotNil)
	// Container is using image, error message should contain errSubstr
	c.Assert(out, checker.Contains, errSubstr, check.Commentf("Container: %q", cleanedContainerID))

	// make sure it didn't delete the busybox name
	images, _ := dockerCmd(c, "images")
	// The name 'busybox' should not have been removed from images
	c.Assert(images, checker.Contains, "busybox")
}

// TestRemoveASingleTagDoesNotRemoveImage (run a container with original tag)
func (s *DockerSuite) TestRmiTag(c *check.C) {
	imagesBefore, _ := dockerCmd(c, "images", "-a")
	dockerCmd(c, "tag", "busybox", "utest:tag1")
	dockerCmd(c, "tag", "busybox", "utest/docker:tag2")
	dockerCmd(c, "tag", "busybox", "utest:5000/docker:tag3")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+3, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))
	}
	dockerCmd(c, "rmi", "utest/docker:tag2")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+2, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))
	}
	dockerCmd(c, "rmi", "utest:5000/docker:tag3")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n")+1, check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))

	}
	dockerCmd(c, "rmi", "utest:tag1")
	{
		imagesAfter, _ := dockerCmd(c, "images", "-a")
		c.Assert(strings.Count(imagesAfter, "\n"), checker.Equals, strings.Count(imagesBefore, "\n"), check.Commentf("before: %q\n\nafter: %q\n", imagesBefore, imagesAfter))

	}
}

// TestRemoveWhenImageIsInUseDoesUntag (when there are multiple tags)
func (s *DockerSuite) TestRmiForceWithExistingContainers(c *check.C) {
	image := "busybox-clone"

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "build", "--no-cache", "-t", image, "-"},
		Stdin: strings.NewReader(`FROM busybox
MAINTAINER foo`),
	}).Assert(c, icmd.Success)

	dockerCmd(c, "run", "--name", "test-force-rmi", image, "/bin/true")

	dockerCmd(c, "rmi", "-f", image)
}

// TestRemoveImageFromMultipleRepositoriesRequiresForce
func (s *DockerSuite) TestRmiForceWithMultipleRepositories(c *check.C) {
	imageName := "rmiimage"
	tag1 := imageName + ":tag1"
	tag2 := imageName + ":tag2"

	buildImageSuccessfully(c, tag1, build.WithDockerfile(`FROM busybox
		MAINTAINER "docker"`))
	dockerCmd(c, "tag", tag1, tag2)

	out, _ := dockerCmd(c, "rmi", "-f", tag2)
	c.Assert(out, checker.Contains, "Untagged: "+tag2)
	c.Assert(out, checker.Not(checker.Contains), "Untagged: "+tag1)

	// Check built image still exists
	images, _ := dockerCmd(c, "images", "-a")
	c.Assert(images, checker.Contains, imageName, check.Commentf("Built image missing %q; Images: %q", imageName, images))
}

// (unit test) TestRemoveImageNotFound
func (s *DockerSuite) TestRmiContainerImageNotFound(c *check.C) {
	// Build 2 images for testing.
	imageNames := []string{"test1", "test2"}
	imageIds := make([]string, 2)
	for i, name := range imageNames {
		dockerfile := fmt.Sprintf("FROM busybox\nMAINTAINER %s\nRUN echo %s\n", name, name)
		buildImageSuccessfully(c, name, build.WithoutCache, build.WithDockerfile(dockerfile))
		id := getIDByName(c, name)
		imageIds[i] = id
	}

	// Create a long-running container.
	runSleepingContainerInImage(c, imageNames[0])

	// Create a stopped container, and then force remove its image.
	dockerCmd(c, "run", imageNames[1], "true")
	dockerCmd(c, "rmi", "-f", imageIds[1])

	// Try to remove the image of the running container and see if it fails as expected.
	out, _, err := dockerCmdWithError("rmi", "-f", imageIds[0])
	// The image of the running container should not be removed.
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, "image is being used by running container", check.Commentf("out: %s", out))
}

// TODO: this should be a unit test
// #13422
func (s *DockerSuite) TestRmiUntagHistoryLayer(c *check.C) {
	image := "tmp1"
	// Build an image for testing.
	dockerfile := `FROM busybox
MAINTAINER foo
RUN echo 0 #layer0
RUN echo 1 #layer1
RUN echo 2 #layer2
`
	buildImageSuccessfully(c, image, build.WithoutCache, build.WithDockerfile(dockerfile))
	out, _ := dockerCmd(c, "history", "-q", image)
	ids := strings.Split(out, "\n")
	idToTag := ids[2]

	// Tag layer0 to "tmp2".
	newTag := "tmp2"
	dockerCmd(c, "tag", idToTag, newTag)
	// Create a container based on "tmp1".
	dockerCmd(c, "run", "-d", image, "true")

	// See if the "tmp2" can be untagged.
	out, _ = dockerCmd(c, "rmi", newTag)
	// Expected 1 untagged entry
	c.Assert(strings.Count(out, "Untagged: "), checker.Equals, 1, check.Commentf("out: %s", out))

	// Now let's add the tag again and create a container based on it.
	dockerCmd(c, "tag", idToTag, newTag)
	out, _ = dockerCmd(c, "run", "-d", newTag, "true")
	cid := strings.TrimSpace(out)

	// At this point we have 2 containers, one based on layer2 and another based on layer0.
	// Try to untag "tmp2" without the -f flag.
	out, _, err := dockerCmdWithError("rmi", newTag)
	// should not be untagged without the -f flag
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, cid[:12])
	c.Assert(out, checker.Contains, "(must force)")

	// Add the -f flag and test again.
	out, _ = dockerCmd(c, "rmi", "-f", newTag)
	// should be allowed to untag with the -f flag
	c.Assert(out, checker.Contains, fmt.Sprintf("Untagged: %s:latest", newTag))
}

// TestRemoveFailsWhenImageIsAParent
func (*DockerSuite) TestRmiParentImageFail(c *check.C) {
	buildImageSuccessfully(c, "test", build.WithDockerfile(`
	FROM busybox
	RUN echo hello`))

	id := inspectField(c, "busybox", "ID")
	out, _, err := dockerCmdWithError("rmi", id)
	c.Assert(err, check.NotNil)
	if !strings.Contains(out, "image has dependent child images") {
		c.Fatalf("rmi should have failed because it's a parent image, got %s", out)
	}
}
