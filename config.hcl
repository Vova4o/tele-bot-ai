telegram_bot_token     = "your_telegram_bot_token"
telegram_channel_id    = 123456789
database_dsn           = "file:news_feed_bot.db?cache=shared&mode=rwc"
fetch_interval         = "10m"
notification_interval  = "1m"
filter_keywords        = ["keyword1", "keyword2", "keyword3"]
openai_key             = "your_openai_key"
openai_prompt          = "Переведи заголовок на русский и создай увлекательное саммари для данной статьи текст на один абзац. Добавляй интересные детали и ключевые аспекты из самой статьи, чтобы привлечь внимание читателя. Включай смайлики, если они органично вписываются в контекст и могут придать тексту дополнительную динамику. Дай ответ на русском языке."
openai_model           = "gemini-pro"