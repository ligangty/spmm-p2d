package main

import (
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func GetKubeClient(kubeconfig string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

var validURLSchemes = []string{"https://", "http://", "tcp://"}

// func NormalizeServerURL(s string) (string, error) {
// 	// normalize scheme
// 	if !hasScheme(s) {
// 		s = validURLSchemes[0] + s
// 	}

// 	addr, err := url.Parse(s)
// 	if err != nil {
// 		return "", fmt.Errorf("Not a valid URL: %v.", err)
// 	}

// 	// normalize host:port
// 	if strings.Contains(addr.Host, ":") {
// 		_, port, err := net.SplitHostPort(addr.Host)
// 		if err != nil {
// 			return "", fmt.Errorf("Not a valid host:port: %v.", err)
// 		}
// 		_, err = strconv.ParseUint(port, 10, 16)
// 		if err != nil {
// 			return "", fmt.Errorf("Not a valid port: %v. Port numbers must be between 0 and 65535.", port)
// 		}
// 	} else {
// 		port := 0
// 		switch addr.Scheme {
// 		case "http":
// 			port = 80
// 		case "https":
// 			port = 443
// 		default:
// 			return "", fmt.Errorf("No port specified.")
// 		}
// 		addr.Host = net.JoinHostPort(addr.Host, strconv.FormatInt(int64(port), 10))
// 	}

// 	// remove trailing slash if that's the only path we have
// 	if addr.Path == "/" {
// 		addr.Path = ""
// 	}

// 	return addr.String(), nil
// }

// func hasScheme(s string) bool {
// 	for _, p := range validURLSchemes {
// 		if strings.HasPrefix(s, p) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func WhoAmI(clientConfig *restclient.Config) (*userv1.User, error) {
// 	client, err := userv1typedclient.NewForConfig(clientConfig)
// 	if err != nil {
// 		return nil, err
// 	}

// 	me, err := client.Users().Get("~", metav1.GetOptions{})

// 	// if we're talking to kube (or likely talking to kube),
// 	if kerrors.IsNotFound(err) || kerrors.IsForbidden(err) {
// 		switch {
// 		case len(clientConfig.BearerToken) > 0:
// 			// convert their token to a hash instead of printing it
// 			h := sha256.New()
// 			h.Write([]byte(clientConfig.BearerToken))
// 			tokenName := fmt.Sprintf("token-%s", base64.RawURLEncoding.EncodeToString(h.Sum(nil)[:9]))
// 			return &userv1.User{ObjectMeta: metav1.ObjectMeta{Name: tokenName}}, nil

// 		case len(clientConfig.Username) > 0:
// 			return &userv1.User{ObjectMeta: metav1.ObjectMeta{Name: clientConfig.Username}}, nil

// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return me, nil
// }
