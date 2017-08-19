## golang-Word-Filter:Use trie tree (datrie lib) for filter

1. Use the trie tree, based on [datrie](https://linux.thai.net/~thep/datrie/datrie.html) lib, which is a famous trie library coded in C. We invoke the c lib via [cgo](https://golang.org/cmd/cgo/). 
2. Use the gin framework for network service, [gin](https://github.com/gin-gonic/gin) is a simple golang framework, you could use other network framerork for your convenience.
3. This filter support Chinese word and other complex script as I have done some pretreatment for the words: split all characters map its bytes to the 0~255, such as "换" will be map to 63 62.


## Quick Start Guide
post with form-data:
1. http://127.0.0.1:8099/filter/upload
file: upload a text file with keys line by line
level: all the keys restrict level in the file, a int value

2. http://127.0.0.1:8099/filter/delete
key: the key you want to delete

3. http://127.0.0.1:8099/filter/search_all
content
Return
```
[
    {
        "head": 0, // The dirty word start position
        "tail": 3, // Dirty word length
        "level": 1 // Dirty level
    },
    {
        "head": 4,
        "tail": 3,
        "level": 1
    }
]
```

### Build the lib yourself #OPTIONAL#
I have already build the libdatrie and put files need into the project and nothing more you need to do. If you want to build the trie tree library yourself, download from the web and follow its ReadMe:
$ LIB_PATH = XXX/lib_path
$ tar zxvf libdatrie-0.2.5.tar.gz
$ cd libdatrie-0.2.5
$ make clean
$ ./configure --prefix=$LIB_PATH
$ make
$ make install
Then, you will see four folder generated in your lib_path, bin include lib share.



## Future Work
1. Using KMP algorithm optimize the content filter operation.
2. Set up a testing data used as benchmark 