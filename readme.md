# tochhc2md

chmファイル内にあるtoc.hhcからmarkdown形式の目次を生成する。  
mdbookの目次ファイル`SUMMARY.md`に使用する目的で作成。mdbookに特化しているわけではない。  

chmファイルは7zip等で展開できる。  
展開した中にある`toc.hhc`を引数に与える。  

`SUMMARY.md`が生成される。

但し、toc.hhcを変換しただけの為、リンク先がhtmlになっている。  
mdbookで使用する場合は、 生成された`SUMMARY.md`を編集して使用すること。  


