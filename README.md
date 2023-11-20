# unique-params
Tries to make URLs produced by tools such as [waybackurls](https://github.com/tomnomnom/waybackurls) and [gau](https://github.com/lc/gau) more appropriate for automatic scanning.

## Installation
```
go install github.com/martinvks/unique-params@latest
```
## Usage
```
cat urls.txt | unique-params > filtered.txt
```
## What does it do?
URLs that have the same host and path are reduced to a single url with all the unique query paramaters.
```
$ cat urls.txt
https://example.com/search?query=computerphile
https://example.com/search?query=quantum+computing
https://example.com/search?utm_source=google
$ cat urls.txt | unique-params
https://example.com/search?query=computerphile&utm_source=google
```
URLs with a numerical or UUID path segment in the same position are reduced to a single url
```
$ cat urls.txt                
https://example.com/articles/1
https://example.com/articles/2
https://example.com/users/59f16da3-a026-4457-8052-6a9e42656415
https://example.com/users/43c291df-0b3d-440a-ba39-7a38c9a213d4
$ cat urls.txt | unique-params
https://example.com/articles/1
https://example.com/users/59f16da3-a026-4457-8052-6a9e42656415
```
