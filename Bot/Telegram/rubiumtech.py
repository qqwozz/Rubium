import json
import telebot
from telebot import types

"""
Почему мы тянем данные из JSON, а не из базы данных?
Нужно использовать готовые эндпоинты из GO API.

"""

TOKEN = "8816074875:AAGLtkH5kc467RQVkRI0ae_IbsGJVmqsV6Y"
bot = telebot.TeleBot(TOKEN)

# Временное хранилище сессий поиска
USER_SEARCHES = {}
# Хранилище ID главного сообщения бота для каждого юзера: {user_id: message_id}
USER_MENU_MSG = {}


def load_portfolio():
    """Считывает базу данных ваших работ из файла portfolio.json."""
    try:
        with open("portfolio.json", "r", encoding="utf-8") as f:
            return json.load(f)
    except FileNotFoundError:
        return []


# === ФУНКЦИИ ДЛЯ СОЗДАНИЯ КНОПОК ===

def get_main_menu_keyboard():
    """Главная панель: Начать поиск и под ней Info."""
    markup = types.InlineKeyboardMarkup()
    markup.row(types.InlineKeyboardButton("🔍 Найти тетрадь", callback_data="go_to_search"))
    markup.row(types.InlineKeyboardButton("ℹ️ Помощь", callback_data="info_stub"))
    return markup


def get_not_found_keyboard():
    """Клавиатура для окна, когда ничего не найдено: Ввести другой тег или Назад."""
    markup = types.InlineKeyboardMarkup()
    markup.row(types.InlineKeyboardButton("⌨️ Ввести другой тег", callback_data="go_to_search"))
    markup.row(types.InlineKeyboardButton("⬅️ Вернуться назад", callback_data="go_to_main_menu"))
    return markup


def get_pagination_keyboard(current_index, total_count):
    """Интерфейс перелистывания работ по схеме с 4 кнопками с поддержкой Web App."""
    markup = types.InlineKeyboardMarkup()

    # Первый ряд: Стрелка влево, Кнопка-ссылка на Web App, Стрелка вправо
    btn_left = types.InlineKeyboardButton("⬅️", callback_data="nav_left")

    # Ссылка на готовую страницу «В разработке» (позже поменяете на свою)
    web_app_url = "https://rubium.tech/"

    btn_web_app = types.InlineKeyboardButton(
        text=f"🌐 Открыть ({current_index + 1}/{total_count})",
        web_app=types.WebAppInfo(url=web_app_url)
    )

    btn_right = types.InlineKeyboardButton("➡️", callback_data="nav_right")
    markup.row(btn_left, btn_web_app, btn_right)

    # Второй ряд: Отмена (возврат в главное меню)
    markup.row(types.InlineKeyboardButton("❌ Отмена", callback_data="go_to_main_menu"))
    return markup


# === ОБРАБОТЧИКИ КОМАНД И ТЕКСТА ===

@bot.message_handler(commands=['start'])
def send_welcome(message):
    """Стартовое окно (отправляется один раз)."""
    user_id = message.from_user.id

    # ... (код удаления)
    try:
        bot.delete_message(message.chat.id, message.message_id)
    except Exception:
        pass

    # Отправляем единственное базовое сообщение
    sent_msg = bot.send_message(
        chat_id=message.chat.id,
        text="Приветствуем! Начните поиск тетрадей по кнопке ниже:",
        reply_markup=get_main_menu_keyboard()
    )
    # Запоминаем ID этого сообщения, чтобы всегда редактировать именно его
    USER_MENU_MSG[user_id] = sent_msg.message_id


@bot.message_handler(func=lambda message: True)
def handle_tag_input(message):
    """Срабатывает, когда пользователь вводит тег текстом."""
    user_id = message.from_user.id
    tag_to_search = message.text.strip().lower()

    # Мгновенно удаляем текст пользователя из чата, чтобы было чисто
    try:
        bot.delete_message(message.chat.id, message.message_id)
    except Exception:
        pass

    if user_id not in USER_MENU_MSG:
        return

    portfolio = load_portfolio()

    # Фильтруем: тег должен совпадать И у работы должен быть статус "public": true
    found_works = [
        work for work in portfolio
        if tag_to_search in [t.lower() for t in work.get("tags", [])]
        if work.get("public") is True
    ]

    # Если ничего не нашли (или все найденные работы приватные)
    if not found_works:
        try:
            bot.edit_message_text(
                chat_id=message.chat.id,
                message_id=USER_MENU_MSG[user_id],
                text=f"❌ По тегу «{tag_to_search}» ничего не найдено среди доступных работ.\n\nВы можете попробовать другой тег или вернуться в меню.",
                reply_markup=get_not_found_keyboard()
            )
        except Exception:
            pass
        return

    # Если публичные работы найдены — сохраняем сессию
    USER_SEARCHES[user_id] = {
        "results": found_works,
        "current_index": 0
    }

    work = found_works[0]
    total = len(found_works)

    text = f"🔍 Найдено доступных тетрадей по тегу: {total}\n\n" \
           f"📂 **{work['title']}**\n" \
           f"📜 Описание: {work['info']}\n" \
           f"⭐ Рейтинг: {work['rating']}\n"\
           f"👤 Автор: {work['user_name']}\n\n" \
           f"🏷️ Теги: {', '.join(work['tags'])}"

    try:
        bot.edit_message_text(
            chat_id=message.chat.id,
            message_id=USER_MENU_MSG[user_id],
            text=text,
            reply_markup=get_pagination_keyboard(0, total),
            parse_mode="Markdown"
        )
    except Exception:
        pass


# === ОБРАБОТЧИКИ НАЖАТИЙ НА КНОПКИ (CALLBACK QUERIES) ===

@bot.callback_query_handler(func=lambda call: True)
def callback_query(call):
    bot.answer_callback_query(call.id)
    chat_id = call.message.chat.id
    message_id = call.message.message_id
    user_id = call.from_user.id

    # На всякий случай обновляем ID главного сообщения при любом клике
    USER_MENU_MSG[user_id] = message_id

    # 1. Возврат в главное меню
    if call.data == "go_to_main_menu":
        bot.edit_message_text(
            chat_id=chat_id, message_id=message_id,
            text="Приветствуем! Начните поиск тетрадей по кнопке ниже:",
            reply_markup=get_main_menu_keyboard()
        )

    # 2. Ожидание ввода тега
    elif call.data == "go_to_search":
        bot.edit_message_text(
            chat_id=chat_id, message_id=message_id,
            text="⌨️ **Напишите нужный тег** прямо в этот чат (например: `математика` или `физика`).\n\n_Ваше сообщение удалится, а это окно покажет результат._",
            parse_mode="Markdown"
        )

    # 3. Заглушка для кнопки INFO
    elif call.data == "info_stub":
        bot.answer_callback_query(
            call.id,
            text="Это меню помощи. Введите предмет или тему текстом, чтобы найти учебные тетради!",
            show_alert=True
        )

    # 4. Перелистывание работ (Влево/Вправо)
    elif call.data in ["nav_left", "nav_right"]:
        if user_id not in USER_SEARCHES or not USER_SEARCHES[user_id]["results"]:
            bot.answer_callback_query(call.id, "Сессия поиска устарела. Нажмите Начать поиск еще раз.", show_alert=True)
            return

        search_data = USER_SEARCHES[user_id]
        works = search_data["results"]
        idx = search_data["current_index"]
        total = len(works)

        if call.data == "nav_left":
            idx = (idx - 1) % total
        else:
            idx = (idx + 1) % total

        search_data["current_index"] = idx

        work = works[idx]
        text = f"🔍 Найдено доступных тетрадей по тегу: {total}\n\n" \
               f"📂 **{work['title']}**\n" \
               f"📜 Описание: {work['info']}\n" \
               f"⭐ Рейтинг: {work['rating']}\n" \
               f"👤 Автор: {work['user_name']}\n\n" \
               f"🏷️ Теги: {', '.join(work['tags'])}"

        bot.edit_message_text(
            chat_id=chat_id, message_id=message_id,
            text=text, reply_markup=get_pagination_keyboard(idx, total),
            parse_mode="Markdown"
        )

    # 5. Кнопка выбора/счетчика работы
    elif call.data == "select_work":
        if user_id in USER_SEARCHES:
            idx = USER_SEARCHES[user_id]["current_index"]
            work = USER_SEARCHES[user_id]["results"][idx]
            bot.answer_callback_query(call.id, text=f"Вы выбрали: {work['title']}", show_alert=True)


if __name__ == "__main__":
    bot.infinity_polling()
