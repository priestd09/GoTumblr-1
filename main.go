package main

import(
)

func main(){
    ParseConfig()
    recorder := CreateRecorders()
    MakeDownloaderWorkers()
    for _, tumblr_source := range Config.TumblrSources{
        td := MakeTumblrDownloader(tumblr_source.Name, tumblr_source.Suffix, tumblr_source.Url, recorder)
        go td.Start()
    }
    kochan_downloader := MakeKochanDownloader(recorder)
    go kochan_downloader.Start()
    a:=make(chan int)
    <-a
}
