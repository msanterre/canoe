# canoe
Canoe is a simple and fast URL shortener.

It uses Redis as a backend store, which allows canoe to scale to a considerable amount of URLs.

You just need to follow this [guide](https://github.com/msanterre/canoe/wiki/Local-setup) to set it up locally.

### Installation
```
git clone https://github.com/msanterre/canoe.git
```

### Usage

Use `fresh` to run Canoe locally: 

```
go get github.com/pilu/fresh
```

Then switch to the project root and run:

```
fresh
```


### TODO
- [x] Provide dockerfile && docker-compose.yml
- [x] Validation
- [ ] Pretty up README
- [x] Ensure slugs don't collide
- [ ] Stats
- [x] JSON support

### Contributions

Contributions are welcome. Fork this repo and submit a pull request and I will review it.

### License
Copyright (c) 2015 Maxime Santerre

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
