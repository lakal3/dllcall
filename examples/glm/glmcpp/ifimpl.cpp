
#define GLM_FORCE_SSE2
#include <glm/mat4x4.hpp>
#include <glm/vec3.hpp>
#include "if.h"

GoError *MultiplyVectors::Multiply() {
	glm::mat4 wm{ Mat[0],Mat[1],Mat[2],Mat[3],Mat[4],Mat[5],Mat[6],Mat[7],
		Mat[8],Mat[9],Mat[10],Mat[11],Mat[12],Mat[13],Mat[14],Mat[15] };
	
	for (int i = 0; i < Vectors.len; i++) {
		auto p = i * 3;
		auto v = Vectors.data[i];
		glm::vec4 vIn{ v[0], v[1], v[2], 1 };
		auto v2 = wm *  vIn;
		Vectors.data[i][0] = v2[0];
		Vectors.data[i][1] = v2[1];
		Vectors.data[i][2] = v2[2];
	}
	return nullptr;
}