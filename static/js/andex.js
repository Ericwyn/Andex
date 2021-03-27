// 访问不需要权限的目录
function gotoDir(path) {
    if (path === null || path === undefined){
        return
    }
    console.log("goto path: " + path)
    if (path === "/"){
        window.location.href= window.location.origin
    } else {
        window.location.search =  "?p=" + path
    }
}

// 访问需要权限的目录
function gotoPermDir(path) {
    if (path === null || path === undefined){
        return
    }

    mdui.dialog({
        title: '管理员登录',
        content: `
                    <div class="">该路径需要输入密码访问</div>
                    <div class="mdui-textfield">
                        <input class="mdui-textfield-input" autocomplete="off" id="path-req-pw" type="password" placeholder="请输入访问密码"/>
                    </div>
                `,
        buttons: [
            {
                text: '取消'
            },
            {
                text: '确认',
                onClick: function(){
                    AJAX({
                        url: "pathPermRequest",
                        method: AJAX_METHOD.POST,
                        json: {
                            path : path,
                            password: document.getElementById("path-req-pw").value
                        },
                        success: function (res){
                            let resJson = JSON.parse(res)
                            console.log(resJson)
                            if (resJson.code !== "1000") {
                                mdui.alert("授权失败");
                            } else {
                                gotoDir(path)
                            }
                        },
                        fail: function (status){
                            mdui.alert("网络错误: " + status);
                        }
                    })
                }
            }
        ]
    });

}

function downFile() {
    let filePath = window.location.search;
    filePath = filePath.replace("/?p=", "")
    filePath = filePath.replace("?p=", "")
    let newUrl = "download?p=" + filePath
    window.open(newUrl,"_blank")
}

const AJAX_METHOD = {
    POST : "POST",
    GET: "GET",
};
/*
新的 AJAX 方法
    ajaxPostData 里面含有一下内容
        url : 参数
        header: header 信息
        method: 方法，字符串 "POST", "GET", "DELETE",或者直接 AJAX_METHOD.POST 这样
        json: 数据体，一个 json 数据体，以 JSON 格式发送请求
        form: 一个二位数组，带上各种参数, 优先级上面，json 的优先级比 form 的优先级要高, 当 json 和 form 都没有的时候，不发送数据，直接请求

            [
                ["userName", "aaa"],
                ["password", "pw12346"]
            ]

        success: function(resp) 一个回调函数，成功之后执行

            回调的参数是 Request 的返回，ResponseText

        fail: function(status) 一个回调函数，xhr.status !== 200 时候执行，

            回调的参数是 status 状态码

        always：function() 一个回调函数，无论 success 或者 fail，都将会被执行
 */
function AJAX(ajaxPostData) {
    let method = ajaxPostData.method;
    let url = ajaxPostData.url;
    if (url === undefined){
        console.log("AJAX 方法没有设置 url");
        return
    }

    if (method === undefined){
        // console.log("AJAX 方法没有设置 method");
        method = AJAX_METHOD.GET;
    }

    let xhr = new XMLHttpRequest();

    // url format
    let urlStart = ""
    if (url.indexOf("https://") === 0) {
        urlStart = "https://"
    } else if (url.indexOf("http://") === 0) {
        urlStart = "http://"
    }
    if (url !== "") {
        url = url.replace(urlStart, "")
    }
    url = url.split("//").join('/')
    url = urlStart + url;

    xhr.open(method, url, true);

    // 存放 header
    if(ajaxPostData.header !== undefined){
        for(let key in ajaxPostData.header){
            xhr.setRequestHeader(key, ajaxPostData.header[key])
        }
    }

    // 存放数据
    if (ajaxPostData.json !== undefined){
        xhr.setRequestHeader('Content-type', 'application/json;charset-UTF-8');
        xhr.send(JSON.stringify(ajaxPostData.json));
    } else if (ajaxPostData.form !== undefined) {
        let params = ajaxPostData.form;
        if (params !== null) {
            if (method === AJAX_METHOD.GET ){
                console.log("非 POST 方法无法提交 FromData 参数");
                return
            }
            let formData = new FormData();
            for (let i = 0; i < params.length; i++) {
                formData.append(params[i][0], params[i][1])
            }
            xhr.send(formData);
        } else {
            xhr.send();
        }
    }else {
        xhr.send();
    }
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                if (ajaxPostData.success !== undefined){
                    ajaxPostData.success(xhr.responseText)
                }
                if (ajaxPostData.always !== undefined){
                    ajaxPostData.always()
                }
            } else {
                if (ajaxPostData.fail !== undefined){
                    ajaxPostData.fail(xhr.status)
                }
                if (ajaxPostData.always !== undefined){
                    ajaxPostData.always()
                }
            }
        }
    }

}
