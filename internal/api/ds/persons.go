package ds

type Person struct {
	Index       int
	Title       string
	Description string
	ImageName   string
	Price       int
}

func GetPipeline() []Person {
	return []Person{
		{0, "Backend", "des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1des1", "backend.jpg", 300000},
		{1, "Frontend", "des2", "frontend.jpg", 150000},
		{2, "Devops", "des3", "devops.jpg", 450000},
		{3, "Chat Bots", "des5", "chat-bots.png", 50000},
		{4, "IOS", "des6", "swift.jpg", 400000},
		{5, "Android", "des7", "kotlin.jpg", 350000},
		{6, "UX/UI Design", "des8", "ux-design.png", 125000},
	}
}
