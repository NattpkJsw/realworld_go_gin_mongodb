package unittest

import "testing"

type testFindSingleArticle struct {
	slug     string
	userID   int
	isErr    bool
	expected string
}

func TestFindSingleArticle(t *testing.T) {
	tests := []testFindSingleArticle{
		{
			slug:     "",
			userID:   1,
			isErr:    true,
			expected: "get articleID failed: sql: no rows in result set",
		},
		{
			slug:     "how-to-train-your-dragon",
			userID:   0,
			isErr:    false,
			expected: `{"article":{"slug":"how-to-train-your-dragon","title":"How to train your dragon","description":"Ever wonder how?","body":"It takes a Jacobian","taglist":["sun","set"],"createdAt":"2024-02-04T13:57:19.098654","updatedAt":"2024-02-04T13:57:19.098654","favorited":false,"favoritesCount":2,"author":{"username":"jake","bio":"I work at statefarm","image":"https://i.stack.imgur.com/xHWG8.jpg","following":false}}}`,
		},
	}

	articleModule := SetupTest().ArticlesModule()
	for _, test := range tests {
		if test.isErr {
			if _, err := articleModule.Usecase().GetSingleArticle(test.slug, test.userID); err.Error() != test.expected {
				t.Errorf("expect: %v, got: %v", test.expected, err.Error())
			}
		} else {
			result, err := articleModule.Usecase().GetSingleArticle(test.slug, test.userID)
			if err != nil {
				t.Errorf("expected: %v, got: %v", nil, err.Error())
			}
			if CompressToJSON(&result) != test.expected {
				t.Errorf("expected: %v, got %v", CompressToJSON(&result), test.expected)
			}
		}
	}
}
