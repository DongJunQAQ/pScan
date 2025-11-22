package scan

import (
	"errors"
	"fmt"
	"sort"
)

var ( //定义两个错误变量
	ErrExists    = errors.New("host already in the list") //主机已在列表
	ErrNotExists = errors.New("host not in the list")     //主机不在列表
)

type HostsList struct { //此结构体代表可执行端口扫描的主机列表
	Hosts []string
}

func (h1 *HostsList) search(host string) (bool, int) { //在主机列表中搜索主机，使用此方法来确保列表中不存在重复的条目
	sort.Strings(h1.Hosts)                        //按字母顺序对主机列表进行升序排序
	i := sort.SearchStrings(h1.Hosts, host)       //搜索主机，在一个已排序（升序）的字符串切片中查找目标字符串并返回目标的索引
	if i < len(h1.Hosts) && h1.Hosts[i] == host { //索引要小于主机列表长度，且使用此索引得到的主机元素与要查找的主机比对，用以验证是否一致
		return true, i //如果一致说明主机已在列表中返回true和元素索引值
	}
	return false, -1
}

func (h1 *HostsList) Add(host string) error { //为主机列表添加主机
	if found, _ := h1.search(host); found { //添加主机前先搜索一下主机是否已存在，存在返回true，不存在返回false
		return fmt.Errorf("%w: %s", ErrExists, host) //返回主机已存在的错误
	}
	h1.Hosts = append(h1.Hosts, host)
	return nil
}

func (h1 *HostsList) Remove(host string) error { //从主机列表中删除主机
	if found, i := h1.search(host); found { //添加主机前先搜索一下主机是否已存在，存在返回true，不存在返回false
		h1.Hosts = append(h1.Hosts[:i], h1.Hosts[i+1:]...)
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotExists, host) //返回主机不存在的错误
}
