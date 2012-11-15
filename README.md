# gcmp

gcmp takes the relative complement of two directories and writes it to a third
directory while ignoring directory structures.

For directories A, B, and C, the following is done

		C = B \ A = { x ∈ B | x ∉ A}

gcmp traverses A with the original set of files, and B with a new set of files.
All files in B which are not found in A are copied to C.

The directory in which a file is found is irrelevant to the comparison. It only
compares the name of a file. Therefore, it might have unwanted effects when
files with the same name are found in different directories.

gcmp is written in Go.

## Installation

	go get github.com/eglimi/gcmp

## Use Case

gcmp could be useful to find images in a large, unordered set that have not yet
been processed.
