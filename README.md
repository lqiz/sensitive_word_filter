1. Use the trie tree, based on the library [datrie| https://linux.thai.net/~thep/datrie/datrie.html]

2. Use the gin framework for network service, [gin|https://github.com/gin-gonic/gin] is a simple framework.


## how wo build the trie tree library:

$ LIB_PATH = XXX/lib_path
$ tar zxvf libdatrie-0.2.5.tar.gz
$ cd libdatrie-0.2.5
$ make clean
$ ./configure --prefix=$LIB_PATH
$ make
$ make install

Then, you will see four folder generated in your lib_path, bin include lib share

## hybrid go with c++ via cgo.



##

split all characters beyonds 0~255 and map to the 0~255, such as "" will be to e5  85, two map

keywords:
tire tree, gin,


##  question
1. Why the if(!pointer) failed
2. no pointer++
3. a[i++] is wrong
4. var err error
   words, err := trie_search_all(tree, content)



## Future Work
1. Using KMP algorithm optimize the content filter operation.
2. set up a benchmark testing data