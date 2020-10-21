package register

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct/html"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	mailvcode   = make(map[string]string)
	Reglock     sync.Mutex
	MailService = EmailData{
		"",
		"",
		"",
		"",
		"CSO2-Server",
		"注册验证码",
		"",
	}
)

func OnTest(w http.ResponseWriter, r *http.Request) {
	log.Println("get a url", r.URL.RequestURI)
}

func OnRegister() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Register server suffered a fault !")
			fmt.Println(err)
			fmt.Println("Fault end!")
		}
	}()
	MailService.SenderMail = Conf.REGEmail
	MailService.SenderCode = Conf.REGPassWord
	MailService.SenderSMTP = Conf.REGSMTPaddr
	http.HandleFunc("/", Register)
	http.HandleFunc("/bg.jpg", OnJpg)
	http.HandleFunc("/test/<filename>", OnTest)
	fmt.Println("Web is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.REGPort)))
	if Conf.EnableMail != 0 {
		fmt.Println("Mail Service is enabled !")
	} else {
		fmt.Println("Mail Service is disabled !")
	}
	err := http.ListenAndServe(":"+strconv.Itoa(int(Conf.REGPort)), nil)
	if err != nil {
		DebugInfo(1, "ListenAndServe:", err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path, err := GetExePath()
	if err != nil {
		DebugInfo(2, err)
		return
	}
	t, err := template.ParseFiles(path + "/CSO2-Server/assert/register.html")
	if err != nil {
		DebugInfo(2, err)
		return
	}
	if strings.Join(r.Form["on_click"], ", ") == "sendmail" &&
		Conf.EnableMail != 0 {
		addrtmp := strings.Join(r.Form["emailaddr"], ", ")
		wth := WebToHtml{Addr: addrtmp}
		if addrtmp == "" {
			wth.Tip = "提示：邮箱不能为空！"
		} else if IsExistsMail([]byte(addrtmp)) {
			wth.Tip = "提示：该邮箱已注册过！"
		} else {
			Vcode := getrand()
			DebugInfo(2, Vcode)
			Reglock.Lock()
			MailService.TargetMail = addrtmp
			MailService.Content = "您的验证码为：" + Vcode + "<br>" + "请勿告诉他人，如非本人操作请忽略本条邮件。请勿回复。"
			Reglock.Unlock()
			if SendEmailTO(&MailService) != nil {
				wth.Tip = "提示：请输入正确的邮箱！"
			} else {
				wth.Tip = "已发送，请在一分钟之内完成注册！"

				Reglock.Lock()
				mailvcode[addrtmp] = Vcode
				Reglock.Unlock()
				go TimeOut(addrtmp)
			}
		}
		t.Execute(w, wth)
	} else if strings.Join(r.Form["on_click"], ", ") == "register" &&
		Conf.EnableMail != 0 {
		addrtmp := strings.Join(r.Form["emailaddr"], ", ")
		usernametmp := strings.Join(r.Form["username"], ", ")
		passwordtmp := strings.Join(r.Form["password"], ", ")
		vercodetmp := strings.Join(r.Form["vercode"], ", ")
		wth := WebToHtml{UserName: usernametmp, Password: passwordtmp, Addr: addrtmp, VerCode: vercodetmp}
		if addrtmp == "" {
			wth.Tip = "提示：邮箱不能为空！"
			t.Execute(w, wth)
			return
		} else if usernametmp == "" {
			wth.Tip = "提示：用户名不能为空！"
			t.Execute(w, wth)
			return
		} else if passwordtmp == "" {
			wth.Tip = "提示：密码不能为空！"
			t.Execute(w, wth)
			return
		} else if vercodetmp == "" {
			wth.Tip = "提示：验证码不能为空！"
			t.Execute(w, wth)
			return
		} else if IsExistsMail([]byte(addrtmp)) {
			wth.Tip = "提示：该邮箱已注册过！"
			t.Execute(w, wth)
			return
		} else if IsExistsUser([]byte(usernametmp)) {
			wth.Tip = "提示：用户名已存在！"
			wth.UserName = ""
			t.Execute(w, wth)
			return
		} else if mailvcode[addrtmp] == vercodetmp {
			u := GetNewUser()
			u.SetUserName([]byte(usernametmp), []byte(usernametmp))
			u.Password = []byte(fmt.Sprintf("%x", md5.Sum([]byte(usernametmp+passwordtmp))))
			u.UserMail = []byte(addrtmp)
			if tf := AddUserToDB(&u); tf != nil {
				wth.Tip = "提示：数据库错误,注册失败！"
				t.Execute(w, wth)
				return
			}
			wth.Tip = "注册成功!"
			t.Execute(w, wth)
		} else {
			wth.Tip = "提示：验证码不正确！"
			t.Execute(w, wth)
		}
	} else if strings.Join(r.Form["on_click"], ", ") == "register" &&
		Conf.EnableMail == 0 {
		usernametmp := strings.Join(r.Form["username"], ", ")
		passwordtmp := strings.Join(r.Form["password"], ", ")
		wth := WebToHtml{UserName: usernametmp, Password: passwordtmp}
		if usernametmp == "" {
			wth.Tip = "提示：用户名不能为空！"
			t.Execute(w, wth)
			return
		} else if passwordtmp == "" {
			wth.Tip = "提示：密码不能为空！"
			t.Execute(w, wth)
			return
		} else if IsExistsUser([]byte(usernametmp)) {
			wth.Tip = "提示：用户名已存在！"
			wth.UserName = ""
			t.Execute(w, wth)
			return
		} else {
			u := GetNewUser()
			u.SetUserName([]byte(usernametmp), []byte(usernametmp))
			u.Password = []byte(fmt.Sprintf("%x", md5.Sum([]byte(usernametmp+passwordtmp))))
			u.UserMail = []byte("Unkown")
			if tf := AddUserToDB(&u); tf != nil {
				wth.Tip = "提示：数据库错误,注册失败！"
				t.Execute(w, wth)
				return
			}
			wth.Tip = "注册成功!"
			t.Execute(w, wth)
		}
	} else {
		t.Execute(w, nil)
	}
}

func OnJpg(w http.ResponseWriter, r *http.Request) {
	path, err := GetExePath()
	if err != nil {
		DebugInfo(2, err)
		return
	}
	file, err := os.Open(path + "/CSO2-Server/assert/bg.jpg")
	if err != nil {
		DebugInfo(2, err)
		return
	}
	buff, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		DebugInfo(2, err)
		return
	}
	w.Write(buff)
}

func getrand() string {
	rand.Seed(time.Now().Unix())
	randnums := strconv.Itoa(rand.Intn(10)) +
		strconv.Itoa(rand.Intn(10)) +
		strconv.Itoa(rand.Intn(10)) +
		strconv.Itoa(rand.Intn(10))
	return randnums
}

func TimeOut(addrtmp string) {
	timer := time.NewTimer(time.Minute)
	<-timer.C

	Reglock.Lock()
	defer Reglock.Unlock()
	delete(mailvcode, addrtmp)
}
