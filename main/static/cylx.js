(function(win){
   
    let cylx = {};
    // request class
    cylx.copy = function(desc,src){
        for(key in src) {
            desc[key] = src[key];
        }
    }
    cylx.HttpRequest = function(){
        this.url = "";
        this.method = "GET";
        this.headers = {};
        this.body = "";
        this.success = function(data,xhq){};
        this.fail = function(e,xhq){};
        this.config = {};
        this.upload = {};
    };
    //http client
    cylx.ajaxBin = function(userReq){
        let oreq = new cylx.HttpRequest();
        cylx.copy(oreq,userReq);
        let cli = new XMLHttpRequest();
        cli.onreadystatechange = function() {
          if(this.readyState == this.DONE) {
              var status = cli.status;
              if (status >= 200 && status < 400) {
                    oreq.success(this.response,cli);
              } else {
                  oreq.fail(this.responseText,cli);
              }
          }
        }
        cli.responseType = "blob";
        //set upload event
        cylx.copy(cli.upload,oreq.upload);
        //handle network error
        cli.onerror= function(e) {
            oreq.fail(e,cli);
        }
        cli.open(oreq.method,oreq.url,true);
        //set headers
        for(k in oreq.headers) {
            let val = oreq.headers[k];
            if(typeof(val) === "string") {
                cli.setRequestHeader(k,oreq.headers[k]);
            }
        }
        //limit method
        cli.send(oreq.body);
    }

    //http client
    cylx.ajax = function(userReq){

        let oreq = new cylx.HttpRequest();
        //copy user's request to default request
        cylx.copy(oreq,userReq);

        let cli = new XMLHttpRequest();
        
        cli.onreadystatechange = function() {
          //done
          if(this.readyState == this.DONE) {
              //console.info("body type",typeof(this.responseText));
              //console.info("body = ",this.responseText);
              //console.info("status = ",cli.status);
              var status = cli.status;
              if (status >= 200 && status < 400) {
                  try{
                    let retData = this.responseText;
                    var contentType = this.getResponseHeader("Content-Type");
                    if(contentType === "application/json") {
                        retData = JSON.parse(this.responseText);
                    }
                    oreq.success(retData,cli);
                  }catch(e){
                    oreq.fail(e,cli);
                  }
              } else {
                  //status code 0 500 400 100
                  oreq.fail(this.responseText,cli);
              }
          }

        }

        //
        cli.onprogress = function(e) {
            let logText = `${e.type}: ${e.loaded} bytes transferred\n`;
            //console.info(logText);
        }

        //set upload event
        cylx.copy(cli.upload,oreq.upload);

        //handle network error
        cli.onerror= function(e) {
            oreq.fail(e,cli);
        }

        cli.open(oreq.method,oreq.url,true);

        //set headers
        for(k in oreq.headers) {
            let val = oreq.headers[k];
            if(typeof(val) === "string") {
                cli.setRequestHeader(k,oreq.headers[k]);
            }
        }
        //limit method
        cli.send(oreq.body);
    };

    cylx.promiseAjax = function(userReq){
        return new Promise(function(resolve, reject){
            let ret = {};
            userReq.success = function(r,xhq){
                ret.error = false
                ret.body = r;
                ret.xhq = xhq;
                resolve(ret);
            };
            userReq.fail = function(e,xhq){
                ret.error = true
                ret.body = e;
                ret.xhq = xhq;
                //throw a exception
                //reject(ret);
                resolve(ret);
            };
            cylx.ajax(userReq);
        });
    }
    cylx.promiseAjaxBin = function(userReq){
        return new Promise(function(resolve, reject){
            let ret = {};
            userReq.success = function(r,xhq){
                ret.error = false
                ret.body = r;
                ret.xhq = xhq;
                resolve(ret);
            };
            userReq.fail = function(e,xhq){
                ret.error = true
                ret.body = e;
                ret.xhq = xhq;
                resolve(ret);
            };
            cylx.ajaxBin(userReq);
        });
    }

    if(!win.cylx) {
        win.cylx = cylx;
    }
    if(win.cylxLibName && typeof(win.cylxLibName)==="string") {
        win[win.cylxLibName] = cylx;
    }
    return cylx;
})(window);
