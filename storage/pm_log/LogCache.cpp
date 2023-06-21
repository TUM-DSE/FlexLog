#include "LogCache.hpp"
#include <iostream>

size_t MAX_CACHE_SIZE = 0;

LogCache::LogCache() {
	logCachePtr = std::make_shared<GSNmapCache>();
	highest_gsn = 0;
}

LogCache::~LogCache() {
	logCachePtr->clear();
}

uint64_t LogCache::next_gsn(uint64_t gsn) {
	GSNmapCache::accessor acc;
	uint64_t next_gsn = gsn + 1;

	if (next_gsn > highest_gsn)
		return gsn;

	while (!(logCachePtr->find(acc, next_gsn)))
		next_gsn++;

	return next_gsn;
}

int LogCache::Append(const std::string record, uint64_t gsn) {
	try {		
		while (logCachePtr->size() >= MAX_CACHE_SIZE) {

			if (logCachePtr->erase(lowest_gsn)) {
				uint64_t tmp = next_gsn(lowest_gsn);
				uint64_t tmp_lowest_gsn = next_gsn(lowest_gsn);

				__atomic_compare_exchange_n(&(lowest_gsn), &tmp_lowest_gsn, tmp, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST);
			}
		}

		if (!(logCachePtr->insert(std::make_pair(gsn, record))))
			return -1;		

		uint64_t tmp_highest_gsn = highest_gsn;
		__atomic_compare_exchange_n(&(highest_gsn), &tmp_highest_gsn, gsn, false, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST);

		return 0;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

const char* LogCache::Read(uint64_t gsn, uint64_t *n_gsn) {
	GSNmapCache::accessor acc;

	try {
		uint64_t ret = next_gsn(gsn);

		if (!(logCachePtr->find(acc, gsn)))
			return "";

		*n_gsn = ret;

		return acc->second.c_str();
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}

int LogCache::Erase(uint64_t gsn) {
	try {
		if (!(logCachePtr->erase(gsn)))
			return 0;
		else
			return 1;
	}
	catch (const std::runtime_error &e){
		std::cerr << e.what();
		exit(-1);
	}	
}
