package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

/*MountType The mount type. Available types:

- `bind` Mounts a file or directory from the host into the container.
  Must exist prior to creating the container.
- `volume` Creates a volume with the given name and options (or uses
  a pre-existing volume with the same name and options). These are
  **not** removed when the container is removed.


swagger:model MountType
*/
type MountType string
