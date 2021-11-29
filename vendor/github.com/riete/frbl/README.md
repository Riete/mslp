tool for continuous reading file by line while file is still writing or is rotated 

```
f := NewFileReader("/path/to/file")
defer f.Close()
go func(f FileReader) {
	for m := range f.Content() {
		fmt.Print(m)
	}
}(f)
for {
	if err := f.ReadLine(); err != nil {
		fmt.Println(err)
		close(f.Content())
		break
	}
	time.Sleep(time.Second)
}
```