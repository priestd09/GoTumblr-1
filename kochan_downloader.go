package main

import(
    "github.com/moovweb/gokogiri"
    "github.com/moovweb/gokogiri/xml"
    "os"
    "io/ioutil"
    "path"
    "time"
    "net/url"
)


type KonachanImg struct{
    id string
    url string
}

func (self *KonachanImg) GetUrl() string{
    return self.url
}

type KochanDownloader struct{
    recorder Recorder
    UrlChan chan ImgResource
    ContenChan chan Content
}

func MakeKochanDownloader(recorder Recorder) KochanDownloader{
    downloader := KochanDownloader{
        recorder:recorder,
        ContenChan:make(chan Content),
        UrlChan:make(chan ImgResource),
    }
    return downloader
}

func (self *KochanDownloader) get_img_id(url string) (string, error){
    return path.Base(url), nil
}

func (self *KochanDownloader) AfterFinished(){
    dir, _ := os.Getwd()
    dir = path.Join(dir, "Downloads", "kochan")
    _, err := os.Stat(dir)
    _, ok := err.(*os.PathError)
    if ok{
        Info("make dir : %s", dir)
        os.MkdirAll(dir, 0775)
    }
    for {
        content := <- self.ContenChan
        if err == nil{
            filename := path.Base(content.Resource.GetUrl())
            filename, _ = url.QueryUnescape(filename)
            filename = path.Join(dir, filename)
            ioutil.WriteFile(filename, content.Content, 0600)
            self.recorder.MarkAsFinished(content.Resource.(*KonachanImg).id)
            Info("New Download Saved : %s", filename)
        }
    }
}

func (self *KochanDownloader) AddUrl(tumblr_img ImgResource){
    self.UrlChan <- tumblr_img
}

func (self *KochanDownloader) ProcessUrl(url string){
    var imgs []xml.Node
    id, err := self.get_img_id(url)
    if err == nil {
        if !self.recorder.HasFinished(id){
            html := Download(url)
            tree, _ := gokogiri.ParseHtml(html)
            png_imgs, _ := tree.Search("//a[@id=\"png\"]/@href")
            jpg_imgs, _ := tree.Search("//a[@id=\"highres\"]/@href")

            if len(png_imgs) > 0{
                imgs = png_imgs
            }else{
                imgs = jpg_imgs
            }
            for _, img_url := range imgs{
                tumblr_img := KonachanImg{
                    id:id,
                    url:img_url.String(),
                }
                self.AddUrl(&tumblr_img)
            }
        }else{
            if Config.Verbos{
                Info("Already Downloaded : %s", id)
            }
        }
    }
}

func (self *KochanDownloader) get_image_list() []string{
    text := Download("http://konachan.com/post/atom")
    tree, _ := gokogiri.ParseXml(text)
    tree.XPathCtx.RegisterNamespace("n", "http://www.w3.org/2005/Atom")
    results, _ := tree.Search("//n:entry/n:link[@rel=\"alternate\"]/@href")
    ret := make([]string, 0)
    for _, r := range results{
        ret = append(ret, r.String())
    }
    return ret
}

func (self *KochanDownloader) check_rss(){
    for {
        image_list := self.get_image_list()
        Info("Getting Rss Konachan: getted %d", len(image_list))
        for _, url := range image_list{
            go self.ProcessUrl(url)
        }
        time.Sleep(time.Minute * Config.CheckInterval)
    }
}

func (self *KochanDownloader) Start(){
    go self.check_rss()
    go self.AfterFinished()
    for {
        img_resource, ok := <- self.UrlChan
        if !ok {
            break
        }
        go Download_raw(img_resource, self)
    }
}

func (self *KochanDownloader) GetContentChan() chan Content{
    return self.ContenChan
}

func (self *KochanDownloader) GetUrlChan() chan ImgResource{
    return self.UrlChan
}
