## packer-builder-oracle-ocisurrogate

 oracle-ocisurrogate is a [Packer plugin](https://www.packer.io/docs/extending/plugins.html) that is able to create new custom images for use with [Oracle Cloud Infrastructure](https://docs.cloud.oracle.com/iaas/Content/GSG/Concepts/concepts.htm) (OCI) that require configuration steps that the standard [oracle-oci](https://www.packer.io/docs/builders/oracle-oci.html) builder is not able to perform..


This builder allows the creation of custom images running different operating systems, custom boot volume file system layouts and technologies and any complex interaction with the operating system that are not possible when running on a mounted root os partition


The builder takes a base image, creates an additional boot volume and allows the execution of steps and the provisioning necessary on the added volume after launching it, and finally snapshots it creating a reusable custom image.

It is recommended that you familiarize yourself with the Key OCI Concepts and Terminology prior to using this builder if you have not done so already.

The builder does not manage images. Once it creates an image, it is up to you to use it or delete it.

This is an advanced builder If you’re just getting started with Packer, I recommend starting with the oracle-oci builder, which is much easier to use.


### Developing packer-builder-oracle-ocisurrogate

#### Packer integration

To run Terragrunt locally, use the `go build` command:

```bash
go build
```
This will create an executable named packer-builder-oracle-ocisurrogate
 that you will need to copy in the same folder where the packer binary is located, see the [Installing Plugins](https://www.packer.io/docs/extending/plugins.html#installing-plugins) section on the packer website to get started

#### Dependencies

* packer-builder-oracle-ocisurrogate
 uses `dep`, a vendor package management tool for golang. See the dep repo for
  [installation instructions](https://github.com/golang/dep).

#### Debug logging

If you set the `PACKER_LOG` environment variable to "true", all logging information from this plugin and packer itself will be printed to
stdout when you run the app.

### License

This code is released under the MIT License. See LICENSE.txt.
