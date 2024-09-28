#include "fibon_if.h"


int64_t fibon(int64_t n) {
	if (n > 2) {
		return fibon(n - 1) + fibon(n - 2);
	}
	return 1;
}

GoError *calcFibonacci::calc() {
	if (n < 1) {
		return new GoError("N must be >= 1");
	}
	result = fibon(n);
	return nullptr;
}

GoError * fastcalcFibonacci::fastCalc() {
	if (n < 1) {
		return new GoError("N must be >= 1");
	}
	*result = fibon(n);
	return nullptr;
}

GoError* calcFibonExtra::calc() {
	if (n < 1) {
		return new GoError("N must be >= 1");
	}
	result = fibon(n);
	return nullptr;
}