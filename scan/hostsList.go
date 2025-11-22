package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

func (h1 *HostsList) Load(hostFile string) error { //导入文件中的主机列表
	f, err := os.Open(hostFile) //以只读模式打开指定路径的文件
	if err != nil {
		if errors.Is(err, os.ErrNotExist) { //如果错误是文件不存在则不执行任何操作
			return nil
		}
		return err //如果是其他原因导致无法打开文件则返回错误
	}
	defer f.Close()                //关闭文件
	scanner := bufio.NewScanner(f) //创建扫描器，扫描器会自动读取f中的数据，其默认以行为单位
	for scanner.Scan() {           //循环读取每行数据，每次调用会读取一行数据
		h1.Hosts = append(h1.Hosts, scanner.Text()) //把当前行的主机名追加到h1.Hosts切片中，scanner.Text(): 获取当前扫描到的行内容（字符串类型，已自动去除末尾的换行符）
	}
	return nil
}

func (h1 *HostsList) Save(hostFile string) error { //将主机列表保存至文件中
	output := ""
	for _, h := range h1.Hosts { //遍历主机列表，将主机列表中的数据变成一行一行字符串的格式
		output += fmt.Sprintln(h) //fmt.Sprintln():将h转为字符串并在末尾添加换行符\n
	}
	return os.WriteFile(hostFile, []byte(output), 0644) //将数据写入文件，[]byte(output)：由于文件存储的本质是字节，因为output是字符串，所以需要通过[]byte()将字符串转为字节切片后存入文件
}
