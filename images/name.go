package images

func FormatImageName(name string, namespace string, tag string) string {
	return registryURL + "/" + namespace + "/" + name + ":" + tag
}
