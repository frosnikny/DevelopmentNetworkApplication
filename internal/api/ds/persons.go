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
		{0, "Backend", "Наша услуга по \"Backend\" для сайтов представляет собой профессиональную разработку и поддержку функциональной части вашего веб-проекта. Мы занимаемся созданием и оптимизацией серверной инфраструктуры, баз данных, а также программированием серверной логики, обеспечивающей взаимодействие между пользователем и сайтом. Наши эксперты в области \"backend\" обладают глубокими знаниями в различных технологиях, таких как языки программирования (например, Python, PHP, Ruby), фреймворки (например, Django, Laravel, Ruby on Rails) и системы управления базами данных (например, MySQL, PostgreSQL, MongoDB). Мы также уделяем особое внимание безопасности и масштабируемости вашего сайта, чтобы обеспечить его эффективную работу даже при высокой нагрузке.", "backend.jpg", 300000},
		{1, "Frontend", "Наша услуга по \"Frontend\" для сайтов предлагает профессиональное создание интерфейса и пользовательского опыта вашего веб-проекта. Мы специализируемся на разработке эстетически привлекательного, интуитивно понятного и отзывчивого дизайна, который привлечет ваших пользователей и обеспечит удобство использования вашего сайта.\nНаши опытные фронтенд-разработчики владеют широким спектром технологий, таких как HTML, CSS, JavaScript и фреймворки, включая React, Angular и Vue.js. Мы активно следим за последними тенденциями в дизайне и пользовательском опыте, чтобы создавать современные и инновационные интерфейсы, соответствующие вашим потребностям и бренду.", "frontend.jpg", 150000},
		{2, "Devops", "Наша услуга по \"DevOps\" для сайтов предлагает комплексный подход к разработке, развертыванию и управлению вашим веб-проектом. Мы объединяем лучшие практики разработки программного обеспечения и операционной деятельности, чтобы обеспечить эффективность, надежность и масштабируемость вашего сайта.\nНаша команда DevOps-специалистов обладает глубокими знаниями и опытом в использовании современных инструментов и технологий, таких как контейнеризация (например, Docker), оркестрация (например, Kubernetes), автоматизация развертывания (например, Ansible, Terraform) и непрерывная интеграция/непрерывное развертывание (CI/CD).", "devops.jpg", 450000},
		{3, "Chat Bots", "Наша услуга по \"Chat Bots\" для сайтов предлагает создание и интеграцию чат-ботов, которые обеспечивают автоматизированное взаимодействие с пользователями вашего веб-проекта. Чат-боты являются эффективным инструментом для улучшения пользовательского опыта, предоставления быстрой поддержки и автоматизации рутиных задач.\nМы разрабатываем чат-ботов, которые могут отвечать на вопросы пользователей, предоставлять информацию о ваших услугах и продуктах, помогать с оформлением заказов, а также выполнять другие задачи, соответствующие вашим потребностям. Наши чат-боты могут быть интегрированы с различными платформами общения, включая веб-сайты, мессенджеры и социальные сети.", "chat-bots.png", 50000},
		{4, "IOS", "Наша услуга по \"iOS (Swift)\" для сайтов предлагает разработку высококачественных мобильных приложений для устройств на iOS, используя язык программирования Swift. Мы специализируемся на создании интуитивно понятных и функциональных приложений, которые оптимизированы под экосистему Apple.\nНаша команда разработчиков iOS обладает глубокими знаниями языка Swift, фреймворков iOS и инструментов разработки, таких как Xcode. Мы следим за последними тенденциями в дизайне и пользовательском опыте iOS-приложений, чтобы создавать современные и привлекательные интерфейсы.", "swift.jpg", 400000},
		{5, "Android", "Наша услуга по \"Android (Kotlin)\" для сайтов предлагает разработку качественных мобильных приложений для устройств на платформе Android, с использованием языка программирования Kotlin. Мы специализируемся на создании инновационных и функциональных приложений, которые оптимизированы под экосистему Android.\nНаша команда разработчиков Android обладает глубокими знаниями языка Kotlin, Android SDK и инструментов разработки, таких как Android Studio. Мы следим за последними трендами в дизайне и пользовательском опыте Android-приложений, чтобы создавать современные и привлекательные интерфейсы.", "kotlin.jpg", 350000},
		{6, "UX/UI Design", "Наша услуга по \"UX/UI Design\" для сайтов предлагает создание привлекательного и функционального пользовательского интерфейса (UI) и удобного пользовательского опыта (UX) для вашего веб-проекта. Мы стремимся создать дизайн, который сочетает визуальную привлекательность с интуитивной навигацией, чтобы обеспечить позитивное впечатление и легкость использования вашего сайта.\nНаша команда UX/UI дизайнеров имеет обширный опыт в разработке интерфейсов, которые оптимизированы для разных устройств и разрешений экранов. Мы уделяем внимание деталям, цветовым схемам, типографике, анимации и другим аспектам дизайна, чтобы создать согласованный и привлекательный образ вашего сайта.", "ux-design.png", 125000},
	}
}
