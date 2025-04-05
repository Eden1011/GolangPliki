package main

import (
	"fmt"
	"strings"
	"time"
)

type FileSystemItem interface {
	Name() string
	Path() string
	Size() int64
	CreatedAt() time.Time
	ModifiedAt() time.Time
}

type Readable interface {
	Read(p []byte) (n int, err error)
}

type Writable interface {
	Write(p []byte) (n int, err error)
}

type Directory interface {
	FileSystemItem
	AddItem(item FileSystemItem) error
	RemoveItem(name string) error
	Items() []FileSystemItem
}

var (
	ErrItemExists       = fmt.Errorf("item already exists")
	ErrItemNotFound     = fmt.Errorf("item not found")
	ErrNotImplemented   = fmt.Errorf("operation not implemented")
	ErrPermissionDenied = fmt.Errorf("permission denied")
	ErrNotDirectory     = fmt.Errorf("not a directory")
	ErrIsDirectory      = fmt.Errorf("is a directory")
)

type BaseItem struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
}

func (b *BaseItem) Name() string {
	return b.name
}

func (b *BaseItem) Path() string {
	return b.path
}

func (b *BaseItem) Size() int64 {
	return b.size
}

func (b *BaseItem) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseItem) ModifiedAt() time.Time {
	return b.modifiedAt
}

func (b *BaseItem) setModifiedAt() {
	b.modifiedAt = time.Now()
}

type Plik struct {
	BaseItem
	content []byte
}

func NewPlik(name, path string) *Plik {
	now := time.Now()
	return &Plik{
		BaseItem: BaseItem{
			name:       name,
			path:       path,
			size:       0,
			createdAt:  now,
			modifiedAt: now,
		},
		content: []byte{},
	}
}

func (p *Plik) Read(b []byte) (n int, err error) {
	if len(p.content) == 0 {
		return 0, nil
	}

	n = copy(b, p.content)
	return n, nil
}

func (p *Plik) Write(b []byte) (n int, err error) {
	p.content = append(p.content, b...)
	p.size = int64(len(p.content))
	p.setModifiedAt()
	return len(b), nil
}

type Katalog struct {
	BaseItem
	items map[string]FileSystemItem
}

func NewKatalog(name, path string) *Katalog {
	now := time.Now()
	return &Katalog{
		BaseItem: BaseItem{
			name:       name,
			path:       path,
			size:       0,
			createdAt:  now,
			modifiedAt: now,
		},
		items: make(map[string]FileSystemItem),
	}
}

func (k *Katalog) AddItem(item FileSystemItem) error {
	if _, exists := k.items[item.Name()]; exists {
		return ErrItemExists
	}

	k.items[item.Name()] = item
	k.setModifiedAt()
	return nil
}

func (k *Katalog) RemoveItem(name string) error {
	if _, exists := k.items[name]; !exists {
		return ErrItemNotFound
	}

	delete(k.items, name)
	k.setModifiedAt()
	return nil
}

func (k *Katalog) Items() []FileSystemItem {
	items := make([]FileSystemItem, 0, len(k.items))
	for _, item := range k.items {
		items = append(items, item)
	}
	return items
}

type SymLink struct {
	BaseItem
	target FileSystemItem
}

func NewSymLink(name, path string, target FileSystemItem) *SymLink {
	now := time.Now()
	return &SymLink{
		BaseItem: BaseItem{
			name:       name,
			path:       path,
			size:       0,
			createdAt:  now,
			modifiedAt: now,
		},
		target: target,
	}
}

func (s *SymLink) Target() FileSystemItem {
	return s.target
}

type PlikDoOdczytu struct {
	BaseItem
	content []byte
}

func NewPlikDoOdczytu(name, path string, content []byte) *PlikDoOdczytu {
	now := time.Now()
	return &PlikDoOdczytu{
		BaseItem: BaseItem{
			name:       name,
			path:       path,
			size:       int64(len(content)),
			createdAt:  now,
			modifiedAt: now,
		},
		content: content,
	}
}

func (p *PlikDoOdczytu) Read(b []byte) (n int, err error) {
	if len(p.content) == 0 {
		return 0, nil
	}

	n = copy(b, p.content)
	return n, nil
}

type VirtualFileSystem struct {
	root *Katalog
}

func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{
		root: NewKatalog("root", "/"),
	}
}

func splitPath(path string) []string {
	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

func (vfs *VirtualFileSystem) getOrCreateDirPath(path string, createIfNotExist bool) (Directory, error) {
	if path == "/" {
		return vfs.root, nil
	}

	parts := splitPath(path)
	var currentDir Directory = vfs.root

	for i, part := range parts {
		currentPath := "/" + strings.Join(parts[:i+1], "/")
		items := currentDir.Items()
		found := false

		for _, item := range items {
			if item.Name() == part {
				if dir, ok := item.(Directory); ok {
					currentDir = dir
					found = true
					break
				} else {
					return nil, ErrNotDirectory
				}
			}
		}

		if !found {
			if !createIfNotExist {
				return nil, ErrItemNotFound
			}

			newDir := NewKatalog(part, currentPath)
			err := currentDir.AddItem(newDir)
			if err != nil {
				return nil, err
			}
			currentDir = newDir
		}
	}

	return currentDir, nil
}

func (vfs *VirtualFileSystem) CreateFile(path string, name string) (FileSystemItem, error) {
	dir, err := vfs.getOrCreateDirPath(path, true)
	if err != nil {
		return nil, err
	}

	filePath := path
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	filePath += name

	file := NewPlik(name, filePath)
	err = dir.AddItem(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (vfs *VirtualFileSystem) CreateReadOnlyFile(path string, name string, content []byte) (FileSystemItem, error) {
	dir, err := vfs.getOrCreateDirPath(path, true)
	if err != nil {
		return nil, err
	}

	filePath := path
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	filePath += name

	file := NewPlikDoOdczytu(name, filePath, content)
	err = dir.AddItem(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (vfs *VirtualFileSystem) CreateDirectory(path string, name string) (Directory, error) {
	dir, err := vfs.getOrCreateDirPath(path, true)
	if err != nil {
		return nil, err
	}

	dirPath := path
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}
	dirPath += name

	newDir := NewKatalog(name, dirPath)
	err = dir.AddItem(newDir)
	if err != nil {
		return nil, err
	}

	return newDir, nil
}

func (vfs *VirtualFileSystem) CreateSymLink(path string, name string, target FileSystemItem) (FileSystemItem, error) {
	dir, err := vfs.getOrCreateDirPath(path, true)
	if err != nil {
		return nil, err
	}

	linkPath := path
	if !strings.HasSuffix(linkPath, "/") {
		linkPath += "/"
	}
	linkPath += name

	symLink := NewSymLink(name, linkPath, target)
	err = dir.AddItem(symLink)
	if err != nil {
		return nil, err
	}

	return symLink, nil
}

func (vfs *VirtualFileSystem) FindItem(path string) (FileSystemItem, error) {
	if path == "/" {
		return vfs.root, nil
	}

	parts := splitPath(path)
	parentPath := "/" + strings.Join(parts[:len(parts)-1], "/")
	name := parts[len(parts)-1]

	dir, err := vfs.getOrCreateDirPath(parentPath, false)
	if err != nil {
		return nil, err
	}

	items := dir.Items()
	for _, item := range items {
		if item.Name() == name {
			return item, nil
		}
	}

	return nil, ErrItemNotFound
}

func (vfs *VirtualFileSystem) DeleteItem(path string) error {
	if path == "/" {
		return ErrPermissionDenied
	}

	parts := splitPath(path)
	parentPath := "/" + strings.Join(parts[:len(parts)-1], "/")
	name := parts[len(parts)-1]

	dir, err := vfs.getOrCreateDirPath(parentPath, false)
	if err != nil {
		return err
	}

	return dir.RemoveItem(name)
}

func main() {
	fs := NewVirtualFileSystem()

	homeDir, err := fs.CreateDirectory("/", "home")
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia katalogu home: %v\n", err)
		return
	}
	fmt.Printf("Utworzono katalog: %s\n", homeDir.Path())

	usersDir, err := fs.CreateDirectory("/home", "users")
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia katalogu users: %v\n", err)
		return
	}
	fmt.Printf("Utworzono katalog: %s\n", usersDir.Path())

	textFile, err := fs.CreateFile("/home/users", "dokument.txt")
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia pliku: %v\n", err)
		return
	}
	fmt.Printf("Utworzono plik: %s\n", textFile.Path())

	writableFile, ok := textFile.(*Plik)
	if !ok {
		fmt.Println("Nie można rzutować na typ Plik")
		return
	}

	content := []byte("To jest przykładowa zawartość pliku tekstowego.")
	_, err = writableFile.Write(content)
	if err != nil {
		fmt.Printf("Błąd podczas zapisu do pliku: %v\n", err)
		return
	}

	readBuffer := make([]byte, 100)
	bytesRead, err := writableFile.Read(readBuffer)
	if err != nil {
		fmt.Printf("Błąd podczas odczytu z pliku: %v\n", err)
		return
	}

	fmt.Printf("Odczytano %d bajtów: %s\n", bytesRead, readBuffer[:bytesRead])

	readOnlyFile, err := fs.CreateReadOnlyFile("/home/users", "readonly.txt", []byte("Ten plik jest tylko do odczytu."))
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia pliku tylko do odczytu: %v\n", err)
		return
	}
	fmt.Printf("Utworzono plik tylko do odczytu: %s\n", readOnlyFile.Path())

	symLink, err := fs.CreateSymLink("/home", "link_do_users", usersDir)
	if err != nil {
		fmt.Printf("Błąd podczas tworzenia dowiązania symbolicznego: %v\n", err)
		return
	}
	fmt.Printf("Utworzono dowiązanie symboliczne: %s\n", symLink.Path())

	rootDir, err := fs.FindItem("/")
	if err != nil {
		fmt.Printf("Błąd podczas pobierania katalogu głównego: %v\n", err)
		return
	}

	rootDirAsDirectory, ok := rootDir.(Directory)
	if !ok {
		fmt.Println("Nie można rzutować na interfejs Directory")
		return
	}

	fmt.Println("\nLista elementów w katalogu głównym:")
	for _, item := range rootDirAsDirectory.Items() {
		fmt.Printf("- %s (ścieżka: %s, rozmiar: %d bajtów)\n", item.Name(), item.Path(), item.Size())
	}

	usersDirItems, err := fs.FindItem("/home/users")
	if err != nil {
		fmt.Printf("Błąd podczas pobierania katalogu users: %v\n", err)
		return
	}

	usersDirAsDirectory, ok := usersDirItems.(Directory)
	if !ok {
		fmt.Println("Nie można rzutować na interfejs Directory")
		return
	}

	fmt.Println("\nLista elementów w katalogu users:")
	for _, item := range usersDirAsDirectory.Items() {
		fmt.Printf("- %s (ścieżka: %s, rozmiar: %d bajtów)\n", item.Name(), item.Path(), item.Size())
	}

	err = fs.DeleteItem("/home/users/dokument.txt")
	if err != nil {
		fmt.Printf("Błąd podczas usuwania pliku: %v\n", err)
		return
	}
	fmt.Println("Usunięto plik dokument.txt")

	fmt.Println("\nLista elementów w katalogu users po usunięciu pliku:")
	for _, item := range usersDirAsDirectory.Items() {
		fmt.Printf("- %s (ścieżka: %s, rozmiar: %d bajtów)\n", item.Name(), item.Path(), item.Size())
	}
}
