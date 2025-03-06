package utils

import (
	"testing"

	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

// TestSplitFile 测试 SplitFile 函数
func TestSplitFile(t *testing.T) {
	testCases := []struct {
		size            int64
		expectedSizes   []int64
		expectedOffsets []int64
		expectedErr     error
	}{
		// size <= 1MB
		{
			size:            512 * 1024, // 512KB
			expectedSizes:   []int64{512 * 1024},
			expectedOffsets: []int64{0},
			expectedErr:     nil,
		},
		// 1MB < size <= 100MB
		{
			size:            50 * MB,
			expectedSizes:   []int64{1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB, 1 * MB},
			expectedOffsets: []int64{0, 1 * MB, 2 * MB, 3 * MB, 4 * MB, 5 * MB, 6 * MB, 7 * MB, 8 * MB, 9 * MB, 10 * MB, 11 * MB, 12 * MB, 13 * MB, 14 * MB, 15 * MB, 16 * MB, 17 * MB, 18 * MB, 19 * MB, 20 * MB, 21 * MB, 22 * MB, 23 * MB, 24 * MB, 25 * MB, 26 * MB, 27 * MB, 28 * MB, 29 * MB, 30 * MB, 31 * MB, 32 * MB, 33 * MB, 34 * MB, 35 * MB, 36 * MB, 37 * MB, 38 * MB, 39 * MB, 40 * MB, 41 * MB, 42 * MB, 43 * MB, 44 * MB, 45 * MB, 46 * MB, 47 * MB, 48 * MB, 49 * MB},
			expectedErr:     nil,
		},
		// 100MB < size <= 400MB
		{
			size:            200 * MB,
			expectedSizes:   []int64{4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB, 4 * MB},
			expectedOffsets: []int64{0, 4 * MB, 8 * MB, 12 * MB, 16 * MB, 20 * MB, 24 * MB, 28 * MB, 32 * MB, 36 * MB, 40 * MB, 44 * MB, 48 * MB, 52 * MB, 56 * MB, 60 * MB, 64 * MB, 68 * MB, 72 * MB, 76 * MB, 80 * MB, 84 * MB, 88 * MB, 92 * MB, 96 * MB, 100 * MB, 104 * MB, 108 * MB, 112 * MB, 116 * MB, 120 * MB, 124 * MB, 128 * MB, 132 * MB, 136 * MB, 140 * MB, 144 * MB, 148 * MB, 152 * MB, 156 * MB, 160 * MB, 164 * MB, 168 * MB, 172 * MB, 176 * MB, 180 * MB, 184 * MB, 188 * MB, 192 * MB, 196 * MB},
			expectedErr:     nil,
		},
		// 400MB < size <= 1GB
		{
			size: 500 * MB,
			expectedSizes: []int64{
				8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB,
				8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB,
				8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB,
				8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 8 * MB, 4 * MB,
			},
			expectedOffsets: []int64{
				0, 8 * MB, 16 * MB, 24 * MB, 32 * MB, 40 * MB, 48 * MB, 56 * MB, 64 * MB, 72 * MB, 80 * MB, 88 * MB, 96 * MB, 104 * MB, 112 * MB, 120 * MB, 128 * MB, 136 * MB, 144 * MB, 152 * MB, 160 * MB, 168 * MB, 176 * MB, 184 * MB, 192 * MB, 200 * MB, 208 * MB, 216 * MB, 224 * MB, 232 * MB, 240 * MB, 248 * MB, 256 * MB, 264 * MB, 272 * MB, 280 * MB, 288 * MB, 296 * MB, 304 * MB, 312 * MB, 320 * MB, 328 * MB, 336 * MB, 344 * MB, 352 * MB, 360 * MB, 368 * MB, 376 * MB, 384 * MB, 392 * MB, 400 * MB, 408 * MB, 416 * MB, 424 * MB, 432 * MB, 440 * MB, 448 * MB, 456 * MB, 464 * MB, 472 * MB, 480 * MB, 488 * MB, 496 * MB,
			},
			expectedErr: nil,
		},
		// 1GB < size <= 4GB
		{
			size: 1*GB + KB,
			expectedSizes: []int64{
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB, 10 * MB,
				10 * MB, 10 * MB, 4*MB + KB,
			},
			expectedOffsets: []int64{
				0, 10 * MB, 20 * MB, 30 * MB, 40 * MB, 50 * MB, 60 * MB, 70 * MB, 80 * MB, 90 * MB,
				100 * MB, 110 * MB, 120 * MB, 130 * MB, 140 * MB, 150 * MB, 160 * MB, 170 * MB, 180 * MB, 190 * MB,
				200 * MB, 210 * MB, 220 * MB, 230 * MB, 240 * MB, 250 * MB, 260 * MB, 270 * MB, 280 * MB, 290 * MB,
				300 * MB, 310 * MB, 320 * MB, 330 * MB, 340 * MB, 350 * MB, 360 * MB, 370 * MB, 380 * MB, 390 * MB,
				400 * MB, 410 * MB, 420 * MB, 430 * MB, 440 * MB, 450 * MB, 460 * MB, 470 * MB, 480 * MB, 490 * MB,
				500 * MB, 510 * MB, 520 * MB, 530 * MB, 540 * MB, 550 * MB, 560 * MB, 570 * MB, 580 * MB, 590 * MB,
				600 * MB, 610 * MB, 620 * MB, 630 * MB, 640 * MB, 650 * MB, 660 * MB, 670 * MB, 680 * MB, 690 * MB,
				700 * MB, 710 * MB, 720 * MB, 730 * MB, 740 * MB, 750 * MB, 760 * MB, 770 * MB, 780 * MB, 790 * MB,
				800 * MB, 810 * MB, 820 * MB, 830 * MB, 840 * MB, 850 * MB, 860 * MB, 870 * MB, 880 * MB, 890 * MB,
				900 * MB, 910 * MB, 920 * MB, 930 * MB, 940 * MB, 950 * MB, 960 * MB, 970 * MB, 980 * MB, 990 * MB,
				1000 * MB, 1010 * MB, 1020 * MB,
			},
			expectedErr: nil,
		},
		// 2.05MB
		{
			size:            2*MB + 50*KB,
			expectedSizes:   []int64{1 * MB, 1 * MB, 50 * KB},
			expectedOffsets: []int64{0, 1 * MB, 2 * MB},
			expectedErr:     nil,
		},
		// size > 4GB
		{
			size:            5 * GB,
			expectedSizes:   nil,
			expectedOffsets: nil,
			expectedErr:     errors.ErrOver4GSize,
		},
	}

	for _, tc := range testCases {
		offsets, sizes, err := SplitFile(tc.size)

		// 检查错误
		if err != nil && err.Error() != tc.expectedErr.Error() {
			t.Errorf("输入大小 %d: 期望错误 %v, 得到错误 %v", tc.size, tc.expectedErr, err)
		}

		// 检查分块大小
		if !equalInt64Slice(sizes, tc.expectedSizes) {
			t.Errorf("输入大小 %d: 分块大小\n期望 %v \n得到 %v", tc.size, tc.expectedSizes, sizes)
		}

		// 检查起始偏移位置
		if !equalInt64Slice(offsets, tc.expectedOffsets) {
			t.Errorf("输入大小 %d: 起始偏移位置\n期望 %v \n得到 %v", tc.size, tc.expectedOffsets, offsets)
		}
	}
}

// equalInt64Slice 检查两个 int64 切片是否相等
func equalInt64Slice(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
