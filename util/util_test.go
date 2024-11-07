package util

import (
	"testing"
	"fmt"
)

func TestIsValidTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		// 基础版本号测试
		{"纯版本号", "1.2.3", true},
		{"无效版本号", "1.2", false},
		{"无效版本号带字母", "1.2.a", false},
		
		// 常见前缀测试
		{"v前缀", "v1.2.3", true},
		{"test前缀", "test1.2.3", true},
		{"prod前缀", "prod1.2.3", true},
		{"prod-前缀", "prod-1.2.3", true},
		{"dev前缀", "dev1.2.3", true},
		{"release前缀", "release1.2.3", true},
		
		// 后缀测试
		{"简单后缀", "1.2.3-alpha", true},
		{"复杂后缀", "1.2.3-beta-rc", true},
		{"带下划线后缀", "1.2.3-beta_rc", true},
		
		// 前后缀组合测试
		{"前后缀组合", "v1.2.3-alpha", true},
		{"prod带后缀", "prod-1.2.3-beta", true},
		{"test带后缀", "test-1.2.3-rc1", true},
		
		// 无效格式测试
		{"前缀数字结尾", "test123-1.2.3", true},
		{"后缀数字开头", "1.2.3-123test", true},
		{"前缀特殊字符", "test#1.2.3", false},
		{"后缀特殊字符", "1.2.3-test#", false},
		{"版本号前有点", "test.1.2.3", false},
		{"多个版本号", "1.2.3.4", false},
		
		// 特殊场景测试
		{"复杂前缀", "my-service-1.2.3", true},
		{"复杂后缀", "1.2.3-rc-alpha", true},
		{"复杂前后缀", "my-service-1.2.3-rc-alpha", true},
		{"下划线前缀", "my_service_1.2.3", true},
		{"下划线后缀", "1.2.3-rc_alpha", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTag(tt.tag)
			if result != tt.expected {
				t.Errorf("IsValidTag(%s) = %v, 期望 %v", tt.tag, result, tt.expected)
			}
		})
	}
}

func TestGetTagParts(t *testing.T) {
	tests := []struct {
		name          string
		tag           string
		wantPrefix    string
		wantVersion   string
		wantSuffix    string
		wantErr       bool
	}{
		// 基础版本号测试
		{
			name:        "纯版本号",
			tag:         "1.2.3",
			wantVersion: "1.2.3",
		},
		
		// 前缀测试
		{
			name:        "v前缀",
			tag:         "v1.2.3",
			wantPrefix:  "v",
			wantVersion: "1.2.3",
		},
		{
			name:        "prod前缀带连字符",
			tag:         "prod-1.2.3",
			wantPrefix:  "prod-",
			wantVersion: "1.2.3",
		},
		{
			name:        "test前缀带连字符",
			tag:         "test-1.2.3",
			wantPrefix:  "test-",
			wantVersion: "1.2.3",
		},
		
		// 后缀测试
		{
			name:        "简单后缀",
			tag:         "1.2.3-alpha",
			wantVersion: "1.2.3",
			wantSuffix:  "-alpha",
		},
		{
			name:        "复杂后缀",
			tag:         "1.2.3-beta-rc",
			wantVersion: "1.2.3",
			wantSuffix:  "-beta-rc",
		},
		
		// 前后缀组合测试
		{
			name:        "前后缀组合",
			tag:         "v-1.2.3-alpha",
			wantPrefix:  "v-",
			wantVersion: "1.2.3",
			wantSuffix:  "-alpha",
		},
		{
			name:        "复杂前后缀",
			tag:         "my-service-1.2.3-rc-alpha",
			wantPrefix:  "my-service-",
			wantVersion: "1.2.3",
			wantSuffix:  "-rc-alpha",
		},
		{
			name:        "带下划线的前后缀",
			tag:         "my_service-1.2.3-rc_alpha",
			wantPrefix:  "my_service-",
			wantVersion: "1.2.3",
			wantSuffix:  "-rc_alpha",
		},
		
		// 错误测试
		{
			name:     "无效标签",
			tag:      "invalid-tag",
			wantErr:  true,
		},
		{
			name:     "前缀数字结尾",
			tag:      "test123-1.2.3",
			wantErr:  false,
			wantPrefix: "test123-",
			wantVersion: "1.2.3",
		},
		
		// 特殊场景测试
		{
			name:        "多个连字符前缀",
			tag:         "my-awesome-service-1.2.3",
			wantPrefix:  "my-awesome-service-",
			wantVersion: "1.2.3",
		},
		{
			name:        "多个连字符后缀",
			tag:         "1.2.3-beta-rc-test",
			wantVersion: "1.2.3",
			wantSuffix:  "-beta-rc-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, version, suffix, err := GetTagParts(tt.tag)
			
			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagParts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil {
				return
			}
			
			// 检查结果
			if prefix != tt.wantPrefix {
				t.Errorf("GetTagParts() prefix = %v, want %v", prefix, tt.wantPrefix)
			}
			if version != tt.wantVersion {
				t.Errorf("GetTagParts() version = %v, want %v", version, tt.wantVersion)
			}
			if suffix != tt.wantSuffix {
				t.Errorf("GetTagParts() suffix = %v, want %v", suffix, tt.wantSuffix)
			}
		})
	}
}

func TestGetLatestTag(t *testing.T) {


	tests := []struct {
		name        string
		mockTags    []string
		prefix      string
		n           int
		suffix      string
		want        string
		wantErr     bool
	}{
		{
			name:     "空标签列表",
			mockTags: []string{},
			prefix:   "v",
			n:        2,
			suffix:   "",
			want:     "v0.0.1",
			wantErr:  false,
		},
		{
			name:     "普通版本号递增",
			mockTags: []string{"v1.2.3", "v1.2.4", "v1.2.5"},
			prefix:   "v",
			n:        2,
			suffix:   "",
			want:     "v1.2.6",
			wantErr:  false,
		},
		{
			name:     "带后缀的版本号(无横线)",
			mockTags: []string{"v1.2.3beta", "v1.2.4beta", "v1.2.5beta"},
			prefix:   "v",
			n:        2,
			suffix:   "beta",
			want:     "v1.2.6beta",
			wantErr:  false,
		},
		{
			name:     "带前缀的版本号(无横线)",
			mockTags: []string{"test1.2.3", "test1.2.4", "test1.2.5"},
			prefix:   "test",
			n:        2,
			suffix:   "",
			want:     "test1.2.6",
			wantErr:  false,
		},
		{
			name:     "带前后缀的版本号(无横线)",
			mockTags: []string{"test1.2.3beta", "test1.2.4beta", "test1.2.5beta"},
			prefix:   "test",
			n:        2,
			suffix:   "beta",
			want:     "test1.2.6beta",
			wantErr:  false,
		},
		{
			name:     "次版本号递增",
			mockTags: []string{"test1.2.5", "test1.3.0", "test1.4.0"},
			prefix:   "test",
			n:        1,
			suffix:   "",
			want:     "test1.5.0",
			wantErr:  false,
		},
		{
			name:     "主版本号递增",
			mockTags: []string{"test1.0.0", "test2.0.0", "test3.0.0"},
			prefix:   "test",
			n:        0,
			suffix:   "",
			want:     "test4.0.0",
			wantErr:  false,
		},
		{
			name:     "复杂前缀",
			mockTags: []string{"release20231.0.0", "release20232.0.0"},
			prefix:   "release2023",
			n:        0,
			suffix:   "",
			want:     "release20233.0.0",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestTag(tt.prefix, tt.n, tt.suffix, tt.mockTags)
			if (err != nil) != tt.wantErr {
				big := ""
				for _, tag := range tt.mockTags {
					if CmpTag(tag, big, tt.prefix, tt.n, tt.suffix) {
						big = tag
					}
				}
				t.Errorf("GetLatestTag() error = %v, wantErr %v, tags: %v, big: %v", err, tt.wantErr, tt.mockTags, big)
				return
			}
			if got != tt.want {
				big := ""
				for _, tag := range tt.mockTags {
					if CmpTag(tag, big, tt.prefix, tt.n, tt.suffix) {
						fmt.Println(tag, ">", big)
						big = tag
					} else {
						fmt.Println(tag, "<", big)
					}
				}
				t.Errorf("GetLatestTag() = %v, want %v, tags: %v, big: %v", got, tt.want, tt.mockTags, big)
			}
		})
	}
} 

