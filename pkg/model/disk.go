package model

import (
	"gorm.io/gorm"
)

// DiskType describes the type of disk image, this can also describe the filesystem contained within
type DiskType int

const (
	// DiskTypeRaw is the simplest DiskType of which nothing extra is known
	DiskTypeRaw DiskType = iota
)

// DiskTransferStrategy describes the strategy used to down- and upload a disk image
type DiskTransferStrategy int

const (
	// DiskTransferStrategyHTTP uses HTTP to transfer the disk image
	DiskTransferStrategyHTTP DiskTransferStrategy = iota
)

// DiskCompressionStrategy the various available disk compression strategies
type DiskCompressionStrategy int

const (
	// DiskCompressionStrategyNone doesn't compress
	DiskCompressionStrategyNone DiskCompressionStrategy = iota
	// DiskCompressionStrategyZSTD compresses disk images with zstd.
	DiskCompressionStrategyZSTD
	// DiskCompressionStrategyGZip uses the standard GZip compression algorithm for disks.
	DiskCompressionStrategyGZip
)

// DiskImage describes a single disk image on the machine
type DiskImage struct {
	gorm.Model `json:"-"`

	DiskType                DiskType
	DiskTransferStrategy    DiskTransferStrategy
	DiskCompressionStrategy DiskCompressionStrategy

	// Location is the place on the booting system, where the disk should be written to.
	// This is usually a /dev device, like /dev/sda or /dev/nvme0n1
	Location string
}

// DiskUUID is the linux by-uuid of a disk
type DiskUUID = string

// ImageUUID is a UUID distinguishing each disk image
type ImageUUID string

// Version stores the version of an ImageModel using an UNIX timestamp
type Version struct {
	gorm.Model `json:"-"`

	Version      int64
	ImageModelID uint
}

// ImageModel defines the database structure for storing the metadata about images
type ImageModel struct {
	// You will see quite a few of these around. They suppress the default values that the ORM creates when it gets
	// cast into JSON.
	gorm.Model `json:"-"`

	// Human identifiable name of this image
	Name string

	// Versions are all possible versions of this image, represented as unix
	// timestamps of their creation. A new version is created whenever a reprovisioning
	// takes place, and this image is replaced.
	Versions []Version

	// ImageUUID is a universally unique identifier for images
	UUID ImageUUID `gorm:"uniqueIndex"`

	// DiskUUID is these disks linux by-uuid
	DiskUUID DiskUUID

	// Foreign key for gorm
	UserModelID uint
}

/* Disk Layout on control_server
This contradicts the actual implementation and other documentation.
TODO: Bring this inline, I prefer this method as well
where 'abc' and 'cdf' are ImageUUIDs

/disks
	/abc
		/1.img
		/2.img
		/3.img
		/4.img
	/cdf
		/1.img
		/2.img
*/