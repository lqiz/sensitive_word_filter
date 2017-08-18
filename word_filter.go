// Copyright 2017 luoruiyi.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

// #cgo LDFLAGS: -L. -ldatrie
// #include "datrie/trie.h"
// #include "datrie/triedefs.h"
// #include "datrie/typedefs.h"
// #include "datrie/alpha-map.h"
import "C"


import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"unsafe"
	"errors"
	"bufio"
	"strconv"
)

func main() {

	var err error
	tree, err = trie_new();
	if err != nil {
		panic(err)
	}


	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	router := gin.Default()    //获得路由实例

	// 注册接口
	router.POST("/filter/upload_keys", upload)
	router.POST("/filter/delete", delete)
	router.POST("/filter/search", search)

	router.Run(":8099")
}

const ok = 1
const key_word_max_len = 1024
var tree *C.Trie

func trie_new() (*C.Trie, error){
	var alpha_map *C.AlphaMap

	alpha_map = C.alpha_map_new()
	if alpha_map == nil {
		return nil, errors.New("alpha_map created failed!")
	}

	if C.alpha_map_add_range(alpha_map, 0x00, 0xff) != 0 {
		C.alpha_map_free(alpha_map)
		return nil, errors.New("alpha_map_add_range failed!")
	}

	trie := C.trie_new(alpha_map)
	C.alpha_map_free(alpha_map)
	if trie == nil {
		return nil, errors.New("trie tree created failed!")
	}

	return trie, nil
}

// AlphaChar in libdatrie is just a uint32
func str_2_alphaChar(str string) ([]C.AlphaChar, int) {
	keyword := []byte(str)

	len := len(keyword)
	alpha_key := make([]C.AlphaChar, len+1)

	for i, v := range keyword {
		alpha_key[i] = C.AlphaChar(v)
	}
	alpha_key[len] = C.TRIE_CHAR_TERM

	return alpha_key, len
}

func trie_store(trie *C.Trie, str string, level int) error {

	alpha_key_sc, len := str_2_alphaChar(str)

	if len > key_word_max_len {
		return errors.New("Beyond the keyword size limit!")
	}

	// it is impossible to pass a slice to c array via cgo.
	var alpha_key [key_word_max_len + 1]C.AlphaChar
	copy(alpha_key[:], alpha_key_sc)

	if int(C.trie_store(trie, (*C.AlphaChar)(unsafe.Pointer(&alpha_key)), C.TrieData(level))) != ok {
		return errors.New("Store failed, unknown reason!")
	}

	return nil
}

func trie_delete(trie *C.Trie, str string) error {

	alpha_key_sc, len := str_2_alphaChar(str)

	if len > key_word_max_len {
		return errors.New("Beyond the keyword size limit!")
	}

	// it is impossible to pass a slice to c array via cgo.
	var alpha_key [key_word_max_len + 1]C.AlphaChar
	copy(alpha_key[:], alpha_key_sc)

	if int(C.trie_delete(trie, (*C.AlphaChar)(unsafe.Pointer(&alpha_key)))) != ok {
		return errors.New("Delete failed, unknown reason!")
	}

	return nil
}

func trie_search(trie *C.Trie, str string) ([]fragment, error) {
	c, len := str_2_alphaChar(str)
	t := C.trie_root(trie)
	// trie tree is empty or the content is empty
	if t == nil || c == nil {
		return nil, errors.New("Tree or filter content is empty!")
	}

	var fg []fragment
	for head := 0; head < len; head++ {
		i := head

		if C.trie_state_is_walkable(t, c[i]) != ok {
			C.trie_state_rewind(t)
			continue
		}

		for i < len && C.trie_state_is_walkable(t, c[i]) == ok {

			if C.trie_state_is_single(t) == ok && C.trie_state_is_walkable(t, C.TRIE_CHAR_TERM) == ok {
				break
			}

			C.trie_state_walk(t, c[i])
			i = i + 1
			if C.trie_state_is_walkable(t, C.TRIE_CHAR_TERM) == ok {
				w := fragment{Head: head, Tail: i - head, Level: int(C.trie_state_get_data(t))}
				fg = append(fg, w)
			}
		}

		C.trie_state_rewind(t)
	}

	C.trie_state_free(t)
	return fg, nil
}

type fragment struct {
	Head  int `json:"head"`
	Tail  int `json:"tail"`
	Level int `json:"level"`
}

// store the tree to file, not used.
func trie_save(trie *C.Trie) {
	save := C.CString("tree_save")
	if C.trie_save(trie, save) == ok {
		fmt.Printf("I am OK")
	}
}

var delete = func(c *gin.Context){
	var err error
	if tree == nil {
		err = errors.New("Tree is empty, perhaps it was init failed!")
	}

	key := c.PostForm("key")
	err = trie_delete(tree, key)

	if err != nil {
		c.JSON(200, "success")
	} else {
		c.JSON(200, err)
	}

}


var search = func(c *gin.Context) {
	if tree == nil {
		c.JSON(200, "Tree is empty, perhaps it was init failed !")
		return
	}

	content := c.PostForm("content")
	fmt.Printf("%v\n", content)

	words, err := trie_search(tree, content)


	if err != nil {
		c.JSON(200, err)
		return
	}

	fmt.Printf("%#v\n", words)
	fmt.Printf("%v\n", words)

	c.JSON(200, words)
}

var upload = func(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	levelStr := c.PostForm("level")

	// we need to save it to disk for a back up, only if the program broken
	c.SaveUploadedFile(file, ".")

	if err != nil || file == nil{
		c.JSON(200, "Upload failed !")
		return
	}

	level, err:= strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(200, "Level should be a number!")
		return
	}

	infile, err:= file.Open();
	defer infile.Close()
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("scanner: %v %v\n", text, level)
		trie_store(tree, text, level)
	}

	if err := scanner.Err(); err != nil {
		c.JSON(200, "err!")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("success! '%s' uploaded!", file.Filename))
}

