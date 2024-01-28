package aws

const (
	/*
		In my experience, you only start hitting the limits when you process over
		50 flashcards at once.

		https://docs.aws.amazon.com/polly/latest/dg/limits.html
		https://docs.aws.amazon.com/translate/latest/dg/what-is-limits.html
	*/
	AWS_THROTTLING_TIMEOUT_SECONDS = 5
)
