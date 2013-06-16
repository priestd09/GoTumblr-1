package main

import(
    dwler "github.com/wooparadog/GoDownload"
)

type Parser interface {
    Start() chan string
    AfterFinished(url string)
}

type Downloader interface{
    Download(url string) []byte
}

type SiteDownloader interface{
    GetContentChan() chan Content
    GetUrlChan() chan ImgResource
}

type ImgResource interface{
    GetUrl() string
}

type Content struct{
    Content []byte
    Resource ImgResource
}

var DownloadWorker chan *Downloader

func MakeDownloaderWorkers() {
    var worker_factory func() Downloader
    if Config.UseProxy {
        worker_factory = ProxyDownloaderFactory
    }else{
        worker_factory = DirectDownloaderFactory
    }
    DownloadWorker = make(chan *Downloader, CONCURENT_DOWNLOADS)
    for i:=0;i<CONCURENT_DOWNLOADS;i++{
        downloader := worker_factory()
        DownloadWorker <- &downloader
    }
}

func ProxyDownloaderFactory() Downloader{
    downloader := dwler.MakePDownloader(Config.Proxy, Config.Timeout)
    return downloader
}

func DirectDownloaderFactory() Downloader{
    downloader := dwler.MakeDirectDownloader(Config.Timeout)
    return downloader
}

func Download_raw(img_resource ImgResource, st SiteDownloader){
    worker := *(<-DownloadWorker)
    defer func(){
        DownloadWorker <- &worker
    }()
    content := worker.Download(img_resource.GetUrl())
    if len(content) > 0{
        result := Content{
            Content:content,
            Resource:img_resource,
        }
        st.GetContentChan() <- result
    }else{
        st.GetUrlChan() <- img_resource
    }
}

func Download(url string) []byte{
    p_downloader := *(<- DownloadWorker)
    defer func(){
        DownloadWorker <- &p_downloader
    }()
    return p_downloader.Download(url)
}
