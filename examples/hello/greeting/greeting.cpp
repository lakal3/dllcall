#include "greeting.h"
#include <iostream>
#include <string>

GoError *greeting::Greet() {
    auto current = this;
	std::string hello;
	while (current != nullptr) {
	    // std::string text(reinterpret_cast<const char *>(current->text.data), current->text.len);
	    // hello.append(text);
	    current->text.append(hello);

	    hello.append("\n");
	    current = current->next;
	}
	std::cout << hello;
	return nullptr;
}
