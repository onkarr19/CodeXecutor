package pkg

import (
	"strings"
	"testing"
)

func TestExecuteAndCleanupContainer(t *testing.T) {
	tests := []struct {
		code     string
		language string
		expected string
		input    string
	}{
		{
			code:     `console.log("Hello, Node.js in Docker!");`,
			language: "JavaScript",
			expected: "Hello, Node.js in Docker!\n",
		},
		{
			code: `#include<stdio.h>
#include<string.h>
int main() {
	char s[100]; scanf("%s", s);
	printf("Hello, C in %s!\n", s);
	return 0;
}`,
			language: "C",
			input:    "Docker",
			expected: "Hello, C in Docker!\n",
		},
		{
			code: `#include<iostream>
using namespace std;
int main() {
	string s; cin>>s;
	cout << "Hello, C++ in " << s << "!";
	return 0;
}`,
			language: "C++",
			expected: "Hello, C++ in Docker!\n",
			input:    "Docker",
		},
		{
			code:     `print("Hello, Python in Docker!")`,
			language: "Python",
			expected: "Hello, Python in Docker!",
			input:    "Docker",
		},
		{
			code: `public class Solution {
				public static void main(String[] args) {
					System.out.println("Hello, Java in Docker!");
				}
			}`,
			language: "Java",
			expected: "Hello, Java in Docker!\n",
		},
		{
			code:     `console.log("Hello, JavaScript in Docker!");`,
			language: "JavaScript",
			expected: "Hello, JavaScript in Docker!\n",
		},
		{
			code:     `fmt.Println("Hello, Go in Docker!")`,
			language: "Go",
			expected: "Hello, Go in Docker!\n",
		},
		{
			code:     `echo "Hello, World!"`,
			language: "UnknownLanguage",                           // Testing an unsupported language
			expected: "language UnknownLanguage is not supported", // Expect an error message
		},
		// Add more test cases for other languages and scenarios
	}

	for _, test := range tests {
		t.Run(test.language, func(t *testing.T) {
			output, err := ExecuteAndCleanupContainer(test.code, test.input, test.language)
			if err != nil {
				if !strings.Contains(err.Error(), test.expected) {
					t.Errorf("Expected Error: %s, Got: %s", test.expected, err.Error())
				}
			}
			if strings.Compare(output, test.expected) == 0 {
				t.Errorf("Expected: %s, Got: %s", test.expected, output)
			}
		})
	}
}
