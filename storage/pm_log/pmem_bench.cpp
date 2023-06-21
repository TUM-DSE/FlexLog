#include "PMLog.h"
#include <iostream>
#include <stdlib.h>
#include <fmt/printf.h>
#include <gflags/gflags.h>
#include <thread>
#include <vector>
#include <sys/time.h>
#include <mutex>
#include <climits>

using GFLAGS_NAMESPACE::ParseCommandLineFlags;
using GFLAGS_NAMESPACE::RegisterFlagValidator;
using GFLAGS_NAMESPACE::SetUsageMessage;

DEFINE_int64(reads, 500, "Percentage of read operations to do.");
DEFINE_int32(threads, 1, "Number of concurrent threads to run.");

DEFINE_int32(record_size, 100, "Size of each record to be appended.");
DEFINE_int32(key_size, 16, "Size of each key.");
DEFINE_int32(nb_ops, 5e6, "Number of total operations to be executed.");
DEFINE_string(distribution_type, "kSequential", "Indexes distribution.");

std::mutex mtx;
struct Times {
	Times(uint64_t s, uint64_t e) : start(s), end(e) {};
	uint64_t start, end;
};
std::vector<Times> times;

#if 0

static enum DistributionType StringToDistributionType(const char* ctype) {
	if (!strcasecmp(ctype, "fixed"))                                        
		return kFixed;
	else if (!strcasecmp(ctype, "uniform"))
		return kUniform;
	else if (!strcasecmp(ctype, "normal"))
		return kNormal;

	fmt::print("[{}] Cannot parse distribution type {}\n", __func__, ctype);
	return kFixed;  // default value
}
#endif

uint64_t NowMicros() {
	struct timeval tv;
	gettimeofday(&tv, nullptr);
	return static_cast<uint64_t>(tv.tv_sec) * 1000000 + tv.tv_usec;
}


class KeyGenerator {
	public:
		KeyGenerator(uint64_t num) : num_(num) {}

		uint64_t Next() {
			return rand() % num_;
		}

	private:
		const uint64_t num_;
};


std::unique_ptr<char[]> random_str(size_t sz) {
	return std::make_unique<char[]>(sz);
}


static void fill_log(PMLog& log, std::unique_ptr<char[]>& dummy_record) {
	// loading phase
	for (size_t i = 0ULL; i < FLAGS_nb_ops/2; i++) {
		cAppend(log, dummy_record.get(), i);
		cCommit(log, i, i);
	}
}

static void thread_func(int&& thread_id) {
	auto record = random_str(FLAGS_record_size);
	PMLog log = startUp_idx(thread_id);
	KeyGenerator gen(FLAGS_nb_ops);

	std::unique_ptr<uint64_t> next_gsn = std::make_unique<uint64_t>();

	fill_log(log, record);
	fmt::print("[{}] half of the log written\n", __func__);
	fmt::print("[{}] {} distribution\n", __func__, FLAGS_distribution_type);

	auto now = NowMicros();
	if (FLAGS_distribution_type == "kSequential") {
		for (int i = 0; i < FLAGS_nb_ops; i++) {
			if (rand()%1000 > FLAGS_reads) {
				cAppend(log, record.get(), i);
				cCommit(log, i, i);
			}
			else
				cRead(log, i, next_gsn.get());

			if (i%15125 == 0)
				fmt::print("[{}] {}\r", __func__, i);
		}
	}
	else if (FLAGS_distribution_type == "kRandom") {
		for (int i = 0; i < FLAGS_nb_ops; i++) {
			auto idx = gen.Next();
			if (rand()%1000 > FLAGS_reads) {
				cAppend(log, record.get(), idx);
				cCommit(log, idx, idx);
			}
			else
				cRead(log, idx, next_gsn.get());

			if (i%15125 == 0)
				fmt::print("[{}] {}\r", __func__, i);
		}
	}
	auto end = NowMicros();

	auto time_diff = end-now;

	{
		std::lock_guard<std::mutex> lock(mtx);
		fmt::print("{} --- {} \n", now, end);
		times.emplace_back(now, end);
	}
	fmt::print("[{}] thread={} finished w/ tp={} ops/s (now={} end={})\n", __func__, thread_id, 1000000*(FLAGS_nb_ops*1.0/(1.0*time_diff)), now, end);
	finalize(log);

	fmt::print("[{}] thread={} terminates\n", __func__, thread_id);
}

int main(int args, char* argv[]) {
	FLAGS_distribution_type = "kSequential";
	ParseCommandLineFlags(&args, &argv, true);
	fmt::print("[{}] reads={}\tthreads={}\trecord_size={}\tkey_size={}\tnb_ops={}\n", __func__, FLAGS_reads, FLAGS_threads, FLAGS_record_size, FLAGS_key_size, FLAGS_nb_ops);
	std::vector<std::thread> threads;
	for (auto i = 0ULL; i < FLAGS_threads; i++)
		threads.emplace_back(thread_func, i);

	for (auto& thread: threads)
		thread.join();

	fmt::print("[{}] finished\n", __func__);

	uint64_t _min = UINT64_MAX;
	uint64_t _max = 1;

	for (auto& item : times) {
		if (_min > item.start) {
			_min = item.start;
		}
		if (_max < item.end) {
			_max = item.end;
		}
	}
	auto total_ops = FLAGS_nb_ops*FLAGS_threads;
	auto time_diff = (_max - _min);
	auto tps = total_ops*1.0/(time_diff*1.0);
	fmt::print("[{}] tp={} ops/s (min={} max={})\n", __func__, (1000000*tps), _min, _max);

}
