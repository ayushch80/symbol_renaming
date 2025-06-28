#include "renamer.h"

#ifdef __cplusplus

#include <iostream>
#include <random>
#include <string>
#include <LIEF/LIEF.hpp>

void SymbolRenamingImpl(const std::string& binPath) {
    auto fat = LIEF::MachO::Parser::parse(binPath);
    if (!fat || fat->size() == 0) {
        std::cerr << "❌ Failed to parse binary: " << binPath << "\n";
        return;
    }

    LIEF::MachO::Binary* binary = fat->at(0);  // Get the first architecture

    uint64_t text_va_start = 0, text_va_end = 0;
    for (const auto& section : binary->sections()) {
        if (section.name() == "__text" &&
            section.segment() &&
            section.segment()->name() == "__TEXT") {
            text_va_start = section.virtual_address();
            text_va_end   = text_va_start + section.size();
            break;
        }
    }

    if (text_va_start == 0) {
        std::cerr << "❌ Could not find __TEXT,__text section.\n";
        return;
    }

    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dist('A', 'z');

    auto random_gibberish = [&](size_t len) {
        std::string s;
        while (s.size() < len) {
            char c = static_cast<char>(dist(gen));
            if (std::isalpha(c)) s.push_back(c);
        }
        return s;
    };

    size_t renamed_count = 0;
    for (auto& sym : binary->symbols()) {
        uint64_t addr = sym.value();
        const std::string& orig = sym.name();

        if (addr == 0 || addr < text_va_start || addr >= text_va_end || orig.size() < 2 || orig == "_main")
            continue;

        std::string new_name = "_" + random_gibberish(orig.size() - 1);
        new_name.resize(orig.size(), '_');

        std::cout << "[+] Renaming " << orig << " → " << new_name << "\n";
        sym.name(new_name);
        renamed_count++;
    }


    std::string outPath = binPath + "_patched";
    binary->write(outPath);
    std::cout << "[+] Renamed " << renamed_count << " symbols\n[+] Output written to: " << outPath << "\n";
}

extern "C" void SymbolRenaming(const char* binPath) {
    if (binPath) {
        SymbolRenamingImpl(binPath);
    }
}
#endif