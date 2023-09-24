#include "Windows.h"
#include <string>
#include <msgboxif.h>

void toWString(GoString str, std::wstring& out) {
	auto wl = MultiByteToWideChar(CP_UTF8, 0, str.data, static_cast<int>(str.len), nullptr, 0);
	out.resize(wl + 1, 0);
	MultiByteToWideChar(CP_UTF8, 0, str.data, static_cast<int>(str.len), &(out.at(0)), wl);
}

GoError* msgBox::show()
{
	std::wstring title;
	std::wstring msg;
	toWString(this->title, title);
	toWString(this->message, msg);
	MessageBoxW(nullptr, msg.c_str(), title.c_str(), MB_OK);
	return nullptr;
}