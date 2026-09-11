package license

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "net/http"
 "os"
 "path/filepath"
 "runtime"
 "time"
)

const Product = "RealmsBrowser"

type Response struct { Status string `json:"status"`; Product string `json:"product"`; Expires string `json:"expires"`; Features []string `json:"features"` }

type Manager struct { Key string; Endpoint string }

func New(key, endpoint string) *Manager { return &Manager{Key:key, Endpoint:endpoint} }

func DeviceID() string { h:=sha256.Sum256([]byte(runtime.GOOS+"|"+runtime.GOARCH+"|"+hostname())); return hex.EncodeToString(h[:]) }
func hostname() string { h,_:=os.Hostname(); return h }

func (m *Manager) Check() (*Response,error) {
 if m.Key=="" { return nil,fmt.Errorf("missing license key") }
 req:=fmt.Sprintf("%s/check?product=%s&key=%s&device=%s",m.Endpoint,Product,m.Key,DeviceID())
 c:=&http.Client{Timeout:10*time.Second}
 r,e:=c.Get(req); if e!=nil{return nil,e}; defer r.Body.Close()
 if r.StatusCode!=200{return nil,fmt.Errorf("license server rejected request")}
 var out Response
 if e=json.NewDecoder(r.Body).Decode(&out);e!=nil{return nil,e}
 if out.Status!="valid"{return nil,fmt.Errorf("invalid license")}
 return &out,nil
}

func CachePath() string { p,_:=os.UserConfigDir(); return filepath.Join(p,"RealmsBrowser","license.json") }
