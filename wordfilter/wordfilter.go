package wordfilter

import (
	"bufio"
	"os"
	"strings"
)

var (
	tree = wordTree{root: &trieNode{ChildMap: make(map[rune]*trieNode, 256), End: false}, replaceChar: '*'}
)

type trieNode struct {
	ChildMap map[rune]*trieNode // 本节点下的所有子节点
	Data     string             // 在最后一个节点保存完整的一个内容
	End      bool               // 标识是否最后一个节点
}

type wordTree struct {
	root        *trieNode
	replaceChar rune
}

// AddWord 添加敏感词到前缀树
func (t *wordTree) AddWord(word string) {
	r := t.root
	chars := []rune(word)
	cnt := len(chars)
	for idx := 0; idx < cnt; idx++ {
		ch := chars[idx]
		r = r.AddChild(ch)
	}
	r.End = true
	r.Data = word
}

/*
替换字符
  - chars 要替换的字符数组
  - begin 开始位置
  - end 结束位置
  - replaceChar 替换字符
*/
func (t *wordTree) replaceRune(chars []rune, begin int, end int, replaceChar rune) {
	for i := begin; i < end; i++ {
		chars[i] = replaceChar
	}
}

/*
判断是否包含敏感词
  - startIndex 开始搜索的位置
  - paramText 要搜索的文本
  - paramReplace 是否替换敏感词
  - paramReplaceChar 替换字符
  - return 返回是否找到敏感词，找到的敏感词的结束位置
*/
func (t *wordTree) FindWord(startIndex int, paramText []rune, paramReplace bool, paramReplaceChar rune) (bool, int) {
	found := false
	endIndex := -1
	for range [1]bool{} {
		r := t.root
		if r == nil {
			break
		}
		cnt := len(paramText)

		if startIndex < 0 || startIndex >= cnt {
			break
		}

		foundStart := -1
		i := startIndex

		for ; i < cnt; i++ {
			ch := paramText[i]
			next := r.FindChild(ch)
			if next == nil {
				// 如果没有下一个节点
				if foundStart == -1 {
					// 如果没有匹配到，则下一个字
					continue
				} else {
					// 如果有开始，如果当前节点是脏字, 则表示找到
					if r.End {
						found = true
						endIndex = i
						break
					} else {
						i = foundStart + 1
						r = t.root
						foundStart = -1
						continue
					}
				}
			} else {
				if foundStart == -1 {
					foundStart = i
				}
				r = next
			}
		}

		if i == cnt && foundStart != -1 {
			if r.End {
				found = true
				endIndex = i
			}
		}

		if found && paramReplace {
			t.replaceRune(paramText, foundStart, endIndex, paramReplaceChar)
		}
	}
	return found, endIndex
}

// AddChild 前缀树添加字节点
func (n *trieNode) AddChild(ch rune) *trieNode {
	if n.ChildMap == nil {
		n.ChildMap = make(map[rune]*trieNode, 10)
	}
	child, ok := n.ChildMap[ch]
	if !ok {
		child = &trieNode{ChildMap: nil, End: false}
		n.ChildMap[ch] = child
		return n.ChildMap[ch]
	} else {
		return child
	}
}

// FindChild 前缀树寻找字节点
func (n *trieNode) FindChild(c rune) *trieNode {
	if n.ChildMap == nil {
		return nil
	}

	if trieNode, ok := n.ChildMap[c]; ok {
		return trieNode
	}
	return nil
}

// LoadWordFile 加载敏感词文件
func loadWordFile(paramWordFile string) (int, error) {
	var retErr error = nil
	var retCode int = 0

	for range [1]int{} {
		f, err := os.Open(paramWordFile)
		if err != nil {
			retErr = err
			retCode = 1
			break
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			tree.AddWord(strings.TrimSpace(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			retErr = err
			retCode = 1
			break
		}
	}
	return retCode, retErr
}

// 初始化敏感词树
func Init(wordFile string) (int, error) {
	retCode, retErr := loadWordFile(wordFile)
	return retCode, retErr
}

// 检查字符串中，是否包含敏感词
func HasSensitiveWord(paramText string) bool {
	cnt := len(paramText)
	if cnt == 0 {
		return false
	}
	textChars := []rune(paramText)
	found, _ := tree.FindWord(0, textChars, false, '*')
	return found
}

// 过滤敏感词，并用*替换
func FilterSensitiveWord(paramText string) string {
	result := ""
	cnt := len(paramText)
	if cnt == 0 {
		return paramText
	}
	textChars := []rune(paramText)

	found := true
	endIndex := 0

	FoundCnt := 0

	for found {
		found, endIndex = tree.FindWord(endIndex, textChars, true, '*')
		if !found {
			break
		} else {
			FoundCnt++
		}

	}

	if FoundCnt > 0 {
		result = string(textChars)
	} else {
		result = paramText
	}
	return result
}
