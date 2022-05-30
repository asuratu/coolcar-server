package wechat

type Service struct {
	AppID     string
	AppSecret string
}

func (s *Service) Resolve(code string) (string, error) {
	return fakeOpenID(s.AppID, s.AppSecret, code)

	/*resp, err := weapp.Login(s.AppID, s.AppSecret, code)
	  if err != nil {
	  	return "", fmt.Errorf("weapp.Login: %v", err)
	  }
	  if err := resp.GetResponseError(); err != nil {
	  	return "", fmt.Errorf("weapp response error: %v", err)
	  }
	  return "", nil*/
}

func fakeOpenID(appId, appSecret, code string) (string, error) {
	//return time.Now().Format("2006-01-02 15:04:05") + " : " + appId + appSecret + code, nil
	return appId + appSecret + code, nil
}
