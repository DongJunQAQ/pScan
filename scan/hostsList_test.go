package scan

import (
	"errors"
	"testing"
)

/*
表驱动测试是一种结构化的测试设计模式，
核心思想是：将测试用例的输入（input）、预期输出（expected output）、测试描述（或其他元信息）
整理成 “测试数据表”（通常是切片 / 数组），
再通过循环遍历表格自动执行所有测试用例，避免重复编写测试代码。
即编写涵盖所测试功能的不同变体的测试用例的常见模式称为表驱动测试，在这种类型的测试中将测试用例定义为匿名struct的slice，
其中包含运行测试所需的数据和预期结果。
*/

/*
t.Errorf：报错不停车，仅记录错误信息，当前测试的后续代码会继续执行，适用于非致命错误
t.Fatalf：报错即停车，记录错误信息后，会终止当前测试的goroutine（注意：不会终止整个测试程序，仅终止当前测试/子测试）并打印测试失败终止点，适用于致命错误
*/

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string //子测试的名称/标题
		host      string //主机
		expectLen int    //预期列表长度
		expectErr error  //预期错误类型
	}{ //定义匿名结构体的切片，里面包含了每条子测试的名称、所需数据、预期结果
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, ErrExists},
	}
	for _, tc := range testCases { //遍历每条子测试，逐个执行子测试
		t.Run(tc.name, func(t *testing.T) { //执行子测试，它会在一个独立的goroutine（协程）中执行，一个子测试失败不会影响其他子测试的执行
			h1 := &HostsList{} //创建被测试的对象：HostsList结构体实例（每次子测试都新建，避免用例间污染）
			//初始化主机列表
			if err := h1.Add("host1"); err != nil { //每次循环前先添加host1主机
				t.Fatal(err) //如果初始化失败则直接终止当前子测试
			}

			err := h1.Add(tc.host) //之后再添加子测试中的主机

			//处理“预期有错误”的子测试（仅适用子测试二）
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("预期失败但实际成功\n") //因为到第二条子测试时会添加失败，此时err变量中必有错误，如果没有则立即终止子测试返回致命错误
				}
				if !errors.Is(err, tc.expectErr) { //如果实际错误与预期错误不符
					t.Errorf("预期错误 %q，实际结果 %q\n", tc.expectErr, err) //则记录测试失败，%q将值格式化为带双引号的字符串
				}
				return //提前返回结束此函数，当前用例的“添加失败场景”校验已完成，无需执行后续“添加成功场景”的校验
			}

			//处理“预期无错误”的子测试（仅适用子测试一）
			if err != nil { //如果实际有错误
				t.Fatalf("预期成功但实际失败 %q\n", err) //则致命错误立即终止子测试
			}
			if len(h1.Hosts) != tc.expectLen { //添加成功后，判断主机列表长度是否符合预期
				t.Errorf("预期列表长度 %d，实际长度 %d\n", tc.expectLen, len(h1.Hosts))
			}
			if h1.Hosts[1] != tc.host { //判断列表中索引为1的元素是否是刚添加的主机
				t.Errorf("预期索引1的主机名为 %q ，实际结果 %q\n", tc.host, h1.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisting", "host1", 1, nil},          //删除存在元素的子测试
		{"RemoveNotFound", "host3", 1, ErrNotExists}, //删除不存在元素的子测试
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h1 := &HostsList{}
			//初始化主机列表
			for _, h := range []string{"host1", "host2"} { //每次循环前先添加两个主机
				if err := h1.Add(h); err != nil {
					t.Fatal(err) //如果初始化失败则直接终止当前子测试
				}
			}

			err := h1.Remove(tc.host) //之后再删除子测试中的主机

			//处理“预期有错误”的子测试（仅适用子测试二）
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("预期失败但实际成功\n") //因为到第二条子测试时会删除失败，此时err变量中必有错误，如果没有则立即终止子测试返回致命错误
				}
				if !errors.Is(err, tc.expectErr) { //如果实际错误与预期错误不符
					t.Errorf("预期错误 %q，实际结果 %q\n", tc.expectErr, err) //则记录测试失败
				}
				return //提前返回结束此函数，当前用例的“删除失败场景”校验已完成，无需执行后续“删除成功场景”的校验
			}

			//处理“预期无错误”的子测试（仅适用子测试一）
			if err != nil { //如果实际有错误
				t.Fatalf("预期成功但实际失败 %q\n", err) //则致命错误立即终止子测试
			}
			if len(h1.Hosts) != tc.expectLen { //删除成功后，判断主机列表长度是否符合预期
				t.Errorf("预期列表长度 %d，实际长度 %d\n", tc.expectLen, len(h1.Hosts))
			}
			if h1.Hosts[0] == tc.host { //判断列表的元素是否删除成功，拿删除后索引为0的元素host2与要删除的元素host1对比
				t.Errorf("删除失败，主机 %q 依旧出现在列表中\n", tc.host) //如果一致则说明删除失败，本次子测试也被标记为失败
			}
		})
	}
}
