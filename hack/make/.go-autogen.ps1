﻿<#
.NOTES
    Author:  @jhowardmsft

    Summary: Windows native version of .go-autogen which generates the
             .go source code for building, and performs resource compilation.

.PARAMETER CommitString
     The commit string. This is calculated externally to this script.

.PARAMETER DockerVersion
     The version such as 17.04.0-dev. This is calculated externally to this script.
#>

param(
    [Parameter(Mandatory=$true)][string]$CommitString,
    [Parameter(Mandatory=$true)][string]$DockerVersion
)

$ErrorActionPreference = "Stop"

try {
    if (Test-Path ".\autogen") {
        Remove-Item ".\autogen" -Recurse -Force | Out-Null
    }

    New-Item -ItemType Directory -Path "autogen\winresources\tmp" | Out-Null
    New-Item -ItemType Directory -Path "autogen\winresources\docker" | Out-Null
    New-Item -ItemType Directory -Path "autogen\winresources\dockerd" | Out-Null
    Copy-Item "hack\make\.resources-windows\resources.go" "autogen\winresources\docker"
    Copy-Item "hack\make\.resources-windows\resources.go" "autogen\winresources\dockerd"

    # Generate a version in the form major,minor,patch,build
    $versionQuad=$DockerVersion -replace "[^0-9.]*" -replace "\.", ","

    # Compile the messages
    windmc hack\make\.resources-windows\event_messages.mc -h autogen\winresources\tmp -r autogen\winresources\tmp
    if ($LASTEXITCODE -ne 0) { Throw "Failed to compile event message resources" }

    # If you really want to understand this madness below, search the Internet for powershell variables after verbatim arguments... Needed to get double-quotes passed through to the compiler options.
    # Generate the .syso files containing all the resources and manifest needed to compile the final docker binaries. Both 32 and 64-bit clients.
    $env:_ag_dockerVersion=$DockerVersion
    $env:_ag_gitCommit=$CommitString

    windres -i hack/make/.resources-windows/docker.rc  -o autogen/winresources/docker/rsrc_amd64.syso  -F pe-x86-64 --use-temp-file -I autogen/winresources/tmp -D DOCKER_VERSION_QUAD=$versionQuad --% -D DOCKER_VERSION=\"%_ag_dockerVersion%\" -D DOCKER_COMMIT=\"%_ag_gitCommit%\"
    if ($LASTEXITCODE -ne 0) { Throw "Failed to compile client 64-bit resources" }

    windres -i hack/make/.resources-windows/docker.rc  -o autogen/winresources/docker/rsrc_386.syso    -F pe-i386   --use-temp-file -I autogen/winresources/tmp -D DOCKER_VERSION_QUAD=$versionQuad --% -D DOCKER_VERSION=\"%_ag_dockerVersion%\" -D DOCKER_COMMIT=\"%_ag_gitCommit%\"
    if ($LASTEXITCODE -ne 0) { Throw "Failed to compile client 32-bit resources" }

    windres -i hack/make/.resources-windows/dockerd.rc -o autogen/winresources/dockerd/rsrc_amd64.syso -F pe-x86-64 --use-temp-file -I autogen/winresources/tmp -D DOCKER_VERSION_QUAD=$versionQuad --% -D DOCKER_VERSION=\"%_ag_dockerVersion%\" -D DOCKER_COMMIT=\"%_ag_gitCommit%\"
    if ($LASTEXITCODE -ne 0) { Throw "Failed to compile daemon resources" }
}
Catch [Exception] {
    # Throw the error onto the caller to display errors. We don't expect this script to be called directly 
    Throw ".go-autogen.ps1 failed with error $_"
}
Finally {
    Remove-Item .\autogen\winresources\tmp -Recurse -Force -ErrorAction SilentlyContinue | Out-Null
    $env:_ag_dockerVersion=""
    $env:_ag_gitCommit=""
}
