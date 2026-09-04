cmake_minimum_required(VERSION 3.20)

project(gweb LANGUAGES NONE)

# ---------------------------------------------------------
# Options
# ---------------------------------------------------------

option(GWEB_BUILD_LIBMEDIA "Build libmedia from source" OFF)
option(GWEB_BUILD_TESTS "Run gweb tests" ON)

# ---------------------------------------------------------
# Paths
# ---------------------------------------------------------

set(GWEB_TARGET gweb)

set(GWEB_SOURCE_DIR
    "${CMAKE_CURRENT_SOURCE_DIR}"
)

set(GWEB_BINARY_DIR
    "${CMAKE_CURRENT_BINARY_DIR}"
)

set(GWEB_DEPLOY_DIR
    "${GWEB_SOURCE_DIR}/deploy"
)

set(GWEB_PYTHONPATH
    "${GWEB_SOURCE_DIR}/pkg"
)

set(GWEB_THIRD_PARTY_DIR
    "${GWEB_SOURCE_DIR}/third_party"
)

message(STATUS "GWEB_TARGET          : ${GWEB_TARGET}")
message(STATUS "GWEB_SOURCE_DIR      : ${GWEB_SOURCE_DIR}")
message(STATUS "GWEB_BINARY_DIR      : ${GWEB_BINARY_DIR}")
message(STATUS "GWEB_DEPLOY_DIR      : ${GWEB_DEPLOY_DIR}")
message(STATUS "GWEB_THIRD_PARTY_DIR : ${GWEB_THIRD_PARTY_DIR}")
message(STATUS "GWEB_PYTHONPATH      : ${GWEB_PYTHONPATH}")

# ---------------------------------------------------------
# Go
# ---------------------------------------------------------

find_program(
    GO_EXECUTABLE
    NAMES go
    REQUIRED
)

execute_process(
    COMMAND
        "${GO_EXECUTABLE}" version

    OUTPUT_VARIABLE
        GO_VERSION

    OUTPUT_STRIP_TRAILING_WHITESPACE
)

message(STATUS "GO_EXECUTABLE     : ${GO_EXECUTABLE}")
message(STATUS "GO_VERSION        : ${GO_VERSION}")

# ---------------------------------------------------------
# Python
# ---------------------------------------------------------

find_package(Python3 REQUIRED COMPONENTS
    Interpreter
    Development
)

set(Python3_CGO_INCLUDE_DIRS "")

foreach(_include_dir IN LISTS Python3_INCLUDE_DIRS)
    string(APPEND
        Python3_CGO_INCLUDE_DIRS
        " -I${_include_dir}"
    )
endforeach()

set(Python3_CGO_LIBRARIES "")

foreach(_python_library IN LISTS Python3_LIBRARIES)
    string(APPEND
        Python3_CGO_LIBRARIES
        " ${_python_library}"
    )
endforeach()

configure_file(
    "${GWEB_SOURCE_DIR}/cmake/pyformat_cgo.go.in"
    "${GWEB_SOURCE_DIR}/pkg/pyformat/cgo.go"
    @ONLY
)

message(STATUS "Python3_EXECUTABLE   : ${Python3_EXECUTABLE}")
message(STATUS "Python3_INCLUDE_DIRS : ${Python3_INCLUDE_DIRS}")
message(STATUS "Python3_LIBRARY_DIRS : ${Python3_LIBRARY_DIRS}")
message(STATUS "Python3_LIBRARIES    : ${Python3_LIBRARIES}")
message(STATUS "Python3_VERSION      : ${Python3_VERSION}")
message(STATUS "Python3_CGO_INCLUDE_DIRS  : ${Python3_CGO_INCLUDE_DIRS}")
message(STATUS "Python3_CGO_LIBRARIES     : ${Python3_CGO_LIBRARIES}")

# ---------------------------------------------------------
# FFmpeg (bundled with libmedia)
# ---------------------------------------------------------

set(FFMPEG_ROOT
    "${GWEB_THIRD_PARTY_DIR}/libmedia/deps/ffmpeg-7.1"
    CACHE PATH "FFmpeg root directory"
)

set(FFMPEG_INCLUDE_DIR
    "${FFMPEG_ROOT}/include"
)

set(FFMPEG_LIBRARY_DIR
    "${FFMPEG_ROOT}/lib"
)

message(STATUS "FFMPEG_ROOT          : ${FFMPEG_ROOT}")
message(STATUS "FFMPEG_INCLUDE_DIR   : ${FFMPEG_INCLUDE_DIR}")
message(STATUS "FFMPEG_LIBRARY_DIR   : ${FFMPEG_LIBRARY_DIR}")

# ---------------------------------------------------------
# Media Library
# ---------------------------------------------------------

set(MEDIA_BUILD_TESTS OFF CACHE BOOL
    "Build media tests"
    FORCE
)

set(LIBMEDIA_ROOT
    "${GWEB_THIRD_PARTY_DIR}/libmedia"
)

set(LIBMEDIA_INCLUDE_DIR
    "${LIBMEDIA_ROOT}/include"
)

if(GWEB_BUILD_LIBMEDIA)

    set(LIBMEDIA_LIBRARY_DIR
        "${GWEB_BINARY_DIR}/libmedia"
    )

    add_subdirectory(
        "${LIBMEDIA_ROOT}"
        "${LIBMEDIA_LIBRARY_DIR}"
    )

else()

    set(LIBMEDIA_LIBRARY_DIR
        "${LIBMEDIA_ROOT}/build"
    )

    add_library(media_static STATIC IMPORTED)

    set_target_properties(media_static PROPERTIES
        IMPORTED_LOCATION
            "${LIBMEDIA_LIBRARY_DIR}/libmedia.a"

        INTERFACE_INCLUDE_DIRECTORIES
            "${LIBMEDIA_INCLUDE_DIR}"
    )

endif()

configure_file(
    "${GWEB_SOURCE_DIR}/cmake/libmedia_cgo.go.in"
    "${GWEB_SOURCE_DIR}/pkg/libmedia/cgo.go"
    @ONLY
)

message(STATUS "LIBMEDIA_ROOT           : ${LIBMEDIA_ROOT}")
message(STATUS "LIBMEDIA_INCLUDE_DIR    : ${LIBMEDIA_INCLUDE_DIR}")
message(STATUS "LIBMEDIA_LIBRARY_DIR    : ${LIBMEDIA_LIBRARY_DIR}")

# ---------------------------------------------------------
# Environment
# ---------------------------------------------------------

if(APPLE)

    set(GLOBAL_ENV
        "PYTHONPATH=${GWEB_PYTHONPATH}"
    )

    set(RUN_ENV
        "CGO_LDFLAGS=-Wl,-rpath,${LIBMEDIA_LIBRARY_DIR}"
    )

elseif(UNIX AND NOT APPLE)

    set(GLOBAL_ENV
        "PYTHONPATH=${GWEB_PYTHONPATH}"
    )

    set(RUN_ENV
        "LD_LIBRARY_PATH=${Python3_LIBRARY_DIRS}:${FFMPEG_LIBRARY_DIR}:${LIBMEDIA_LIBRARY_DIR}"
    )

endif()

message(STATUS "GLOBAL_ENV         : ${GLOBAL_ENV}")
message(STATUS "RUN_ENV            : ${RUN_ENV}")

# ---------------------------------------------------------
# Python compileall
# ---------------------------------------------------------

add_custom_target(compile_python
    COMMAND
        "${Python3_EXECUTABLE}"
        -m compileall
        -b
        "${GWEB_SOURCE_DIR}/pkg/pyformat"

    WORKING_DIRECTORY
        "${GWEB_SOURCE_DIR}"

    COMMENT
        "Compile Python files"
)

add_custom_target(clean_python
    COMMAND
        find "${GWEB_SOURCE_DIR}/pkg"
        -name "*.pyc"
        -type f
        -delete

    COMMAND
        find "${GWEB_SOURCE_DIR}/pkg"
        -type d
        -name "__pycache__"
        -exec rm -r "{}" "+"

    COMMAND
        find "${GWEB_SOURCE_DIR}/tests"
        -type d
        -name "__pycache__"
        -exec rm -r "{}" "+"

    COMMENT "Clean python caches"
)

# ---------------------------------------------------------
# FormatUI assets
# ---------------------------------------------------------

add_custom_target(copy_formatui
    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/img"

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/js"

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/css"

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/plugins"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/formatui/dist/img/formatui.svg"
        "${GWEB_SOURCE_DIR}/public/img/format.svg"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/formatui/dist/index.min.js"
        "${GWEB_SOURCE_DIR}/public/js/format.min.js"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/formatui/dist/index.min.css"
        "${GWEB_SOURCE_DIR}/public/css/format.min.css"

    COMMAND
        "${CMAKE_COMMAND}" -E copy_directory
        "${GWEB_SOURCE_DIR}/third_party/formatui/dist/plugins"
        "${GWEB_SOURCE_DIR}/public/plugins"

    COMMENT
        "Copy formatui assets"
)

add_custom_target(clean_formatui
    COMMAND
        "${CMAKE_COMMAND}" -E rm -f
        "${GWEB_SOURCE_DIR}/public/img/format.svg"
        "${GWEB_SOURCE_DIR}/public/js/format.min.js"
        "${GWEB_SOURCE_DIR}/public/css/format.min.css"

    COMMAND
        "${CMAKE_COMMAND}" -E rm -rf
        "${GWEB_SOURCE_DIR}/public/plugins/bootstrap-icons@*"
        "${GWEB_SOURCE_DIR}/public/plugins/clipboard@*"
        "${GWEB_SOURCE_DIR}/public/plugins/json5@*"
        "${GWEB_SOURCE_DIR}/public/plugins/w2ui@*"

    COMMENT "Delete format UI files"
)

# ---------------------------------------------------------
# MediaUI assets
# ---------------------------------------------------------

add_custom_target(copy_mediaui
    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/img"

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/js"

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_SOURCE_DIR}/public/css"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/mediaui/dist/mediaui.svg"
        "${GWEB_SOURCE_DIR}/public/img/mediaui.svg"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/mediaui/dist/mediaui.js"
        "${GWEB_SOURCE_DIR}/public/js/mediaui.js"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_SOURCE_DIR}/third_party/mediaui/dist/mediaui.css"
        "${GWEB_SOURCE_DIR}/public/css/mediaui.css"

    COMMENT
        "Copy mediaui assets"
)

add_custom_target(clean_mediaui
    COMMAND
        "${CMAKE_COMMAND}" -E rm -f
        "${GWEB_SOURCE_DIR}/public/img/mediaui.svg"
        "${GWEB_SOURCE_DIR}/public/js/mediaui.js"
        "${GWEB_SOURCE_DIR}/public/css/mediaui.css"

    COMMENT "Delete media UI files"
)

# ---------------------------------------------------------
# GWeb build
# ---------------------------------------------------------

add_custom_target(${GWEB_TARGET} ALL
    COMMAND
        ${CMAKE_COMMAND} -E env
        ${GLOBAL_ENV}
        "${GO_EXECUTABLE}"
        build
        -v
        -o
        "${GWEB_BINARY_DIR}/${GWEB_TARGET}"
        "cmd/gweb.go"

    WORKING_DIRECTORY
        "${GWEB_SOURCE_DIR}"

    DEPENDS
        media_static
        compile_python
        copy_formatui
        copy_mediaui

    COMMENT
        "Build gweb"
)

# ---------------------------------------------------------
# GWeb deploy
# ---------------------------------------------------------

add_custom_target(deploy_gweb
    DEPENDS
        ${GWEB_TARGET}

    COMMAND
        "${CMAKE_COMMAND}" -E make_directory
        "${GWEB_DEPLOY_DIR}/lib"

    COMMAND
        "${CMAKE_COMMAND}" -E copy
        "${GWEB_BINARY_DIR}/${GWEB_TARGET}"
        "${GWEB_DEPLOY_DIR}/${GWEB_TARGET}"

    COMMAND
        rsync
        -av
        "--include=*/"
        "--include=*.pyc"
        "--exclude=*"
        "${GWEB_SOURCE_DIR}/pkg/pyformat"
        "${GWEB_DEPLOY_DIR}"

    COMMENT
        "Prepare gweb deploy directory"
)

# MacOS configuration
if(APPLE)

    if(NOT GWEB_BUILD_LIBMEDIA)

        find_program(
            INSTALL_NAME_TOOL
            NAMES install_name_tool
            REQUIRED
        )

        add_custom_command(
            TARGET
                deploy_gweb

            POST_BUILD

            COMMAND
                "${CMAKE_COMMAND}" -E copy
                "${LIBMEDIA_LIBRARY_DIR}/libmedia.dylib"
                "${GWEB_DEPLOY_DIR}/lib/libmedia.dylib"

            COMMAND
                "${INSTALL_NAME_TOOL}"
                -add_rpath
                "@executable_path/lib"
                "${GWEB_DEPLOY_DIR}/${GWEB_TARGET}"

            COMMENT
                "Add RPATH to gweb when deploy"
        )

    endif()

# Linux configuration
elseif(UNIX AND NOT APPLE)

    find_program(
        PATCHELF_EXECUTABLE
        NAMES patchelf
    )

    if(PATCHELF_EXECUTABLE)

        add_custom_command(
            TARGET
                ${GWEB_TARGET}

            POST_BUILD

            COMMAND
                "${PATCHELF_EXECUTABLE}"
                --set-rpath
                "\$ORIGIN"
                "${GWEB_BINARY_DIR}/${GWEB_TARGET}"

            COMMENT
                "Set gweb RPATH to \$ORIGIN"
        )

    else()

        message(WARNING
            "patchelf not found; gweb RPATH will not be modified"
        )

    endif()

endif()

add_custom_target(clean_deploy
    COMMAND
        find "${GWEB_SOURCE_DIR}/deploy"
        -mindepth 1
        -name "gweb.yaml"
        -prune
        -o
        -exec rm -rf "{}" "+"

    COMMENT "Clean deploy files"
)

# ---------------------------------------------------------
# GWeb run
# ---------------------------------------------------------

add_custom_target(run_gweb
    COMMAND
        ${CMAKE_COMMAND} -E env
        ${GLOBAL_ENV}
        ${RUN_ENV}
        "${GO_EXECUTABLE}"
        run "cmd/gweb.go"
        -p 5015
        --debug

    WORKING_DIRECTORY
        "${GWEB_SOURCE_DIR}"

    DEPENDS
        ${GWEB_TARGET}

    COMMENT
        "Run gweb"
)

# ---------------------------------------------------------
# GWeb clean
# ---------------------------------------------------------

add_custom_target(clean_gweb
    DEPENDS
        clean_python
        clean_formatui
        clean_mediaui
        clean_deploy

    COMMAND
        "${CMAKE_COMMAND}" -E rm -rf
        "${GWEB_SOURCE_DIR}/gweb"

    COMMENT "Clean project"
)

# ---------------------------------------------------------
# GWeb tests
# ---------------------------------------------------------

if(GWEB_BUILD_TESTS)

    add_custom_target(test_pyformat
        COMMAND
            "${CMAKE_COMMAND}" -E env
            "${GLOBAL_ENV}"
            "${Python3_EXECUTABLE}"
            -m unittest
            -v
            tests/pyformat/test_pytext.py

        COMMAND
            "${CMAKE_COMMAND}" -E env
            "${GLOBAL_ENV}"
            "${Python3_EXECUTABLE}"
            -m unittest
            -v
            tests/pyformat/test_pycmd.py

        COMMAND
            "${CMAKE_COMMAND}" -E env
            "${GLOBAL_ENV}"
            "${Python3_EXECUTABLE}"
            -m unittest
            -v
            tests/pyformat/test_pyfmt.py

        WORKING_DIRECTORY
            "${GWEB_SOURCE_DIR}"

        DEPENDS
            compile_python

        COMMENT
            "Run pyformat tests"
    )

endif()
