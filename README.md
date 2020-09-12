# ImgSearchEngine

## Description
This repository contains the implementation of a image-based search engine using the normalized Hamming distance of the Perception Hashes of images to perform KNN and Range Search. 

## VP-Tree
The data structure used to store the pHashes is the Vantage-Point Tree. 

## Example
An example for serving the search engine can be found inside `src/example`. The 
application will load a tab separated file called `load_file_phash.csv` (not provided) containing 
two columns: `path` and `phash`, where the path refers to the path of the image, and the phash
refers to its Perception Hash. It also provides a File System to store uploaded images under the 
relative directory `images/temp/`. 

To begin serving the example engine, 

```
cd src/example
go run .
```

A simple UI is implemented in React JS under `app/`. The UI lets you upload an image, and search 
similar image against the engine using KNN or Range Search. To start the React app, 

```
cd app
npm start
```
