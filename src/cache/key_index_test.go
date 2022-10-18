package cache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRoot_Root(t *testing.T) {
	var r *Root
	r = r.Root("/a/b/c").
		Root("/a/b/a").
		Root("/b/").
		Root("/a")

	t.Run("Checking tree building", func(t *testing.T) {
		if r.childs[0].name != "a" {
			t.Fatal("should be a")
		}
		if r.childs[1].name != "b" {
			t.Fatal("should be b")
		}
		if r.childs[0].childs[0].name != "b" {
			t.Fatal("should be b")
		}
		if r.childs[0].childs[0].childs[0].name != "a" {
			t.Fatal("should be c")
		}
		if r.childs[0].childs[0].childs[1].name != "c" {
			t.Fatal("should be c")
		}
	})

	t.Run("Get all paths", func(t *testing.T) {
		paths := r.GetAllPaths("")
		expected_paths := []string{
			"/a",
			"/a/b/a",
			"/a/b/c",
			"/b",
		}
		if !reflect.DeepEqual(paths, expected_paths) {
			fmt.Println(paths)
			fmt.Println(expected_paths)
			t.Fatal("Not the same paths")
		}
	})

	t.Run("Get all path of /a", func(t *testing.T) {
		paths := r.GetAllPaths("/a")
		expected_paths := []string{
			"/a",
			"/a/b/a",
			"/a/b/c",
		}
		if !reflect.DeepEqual(paths, expected_paths) {
			fmt.Println(paths)
			fmt.Println(expected_paths)
			t.Fatal("Not the same paths")
		}
	})

	t.Run("Get all path of /a/b", func(t *testing.T) {
		paths := r.GetAllPaths("/a/b")
		expected_paths := []string{
			"/a/b/a",
			"/a/b/c",
		}
		if !reflect.DeepEqual(paths, expected_paths) {
			fmt.Println(paths)
			fmt.Println(expected_paths)
			t.Fatal("Not the same paths")
		}
	})
}
