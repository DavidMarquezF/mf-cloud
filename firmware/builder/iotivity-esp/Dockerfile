FROM ubuntu:20.04 as iotivity-esp

ARG root_directory=/mf-builder
RUN apt-get update
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y git wget flex bison gperf python3 python3-pip python3-setuptools \
 python3-serial python3-click python3-cryptography python3-future python3-pyparsing \
 python3-pyelftools cmake ninja-build ccache libffi-dev libssl-dev libusb-1.0-0 \
 && apt-get autoremove -y \
   && rm -rf /var/lib/apt/lists/* \
   && update-alternatives --install /usr/bin/python python /usr/bin/python3 10

RUN python -m pip install --upgrade pip virtualenv


WORKDIR ${root_directory}

# Clone iotivity
RUN git clone --recursive --branch mf-plgd.dev --single-branch https://github.com/DavidMarquezF/iotivity-lite.git


# Install idf
ARG idf_path=${root_directory}/iotivity-lite/port/esp32/esp-idf
ARG idf_tools_path=${root_directory}/iotivity-lite/port/esp32


ENV IDF_PATH=${idf_path}
ENV IDF_TOOLS_PATH=${idf_tools_path}
WORKDIR ${idf_tools_path}
RUN git clone --recursive https://github.com/espressif/esp-idf.git 

# Submodule init asks for account, give a fake one
RUN git config --global user.email "mf_cont@mf.ca"
RUN git config --global user.name "Mf Container"
# Run patches
RUN (cd esp-idf && git checkout 457ce080aec9811faee26a1ea5293fff17049fee && git submodule init && git submodule update)


# Install all the required tools, plus CMake
RUN $IDF_PATH/tools/idf_tools.py --non-interactive install required \
  && $IDF_PATH/tools/idf_tools.py --non-interactive install cmake \
  && $IDF_PATH/tools/idf_tools.py --non-interactive install-python-env \
  && rm -rf $IDF_TOOLS_PATH/dist


# Ccache is installed, enable it by default
ENV IDF_CCACHE_ENABLE=1

RUN ( cd esp-idf/components/mbedtls/mbedtls && git am ../../../../patches/mbedtls/*.patch )
RUN ( cd esp-idf && find ../patches/esp-idf/ -type f -name '*.patch' -exec patch -p1 -i {} \; )
RUN ( cd esp-idf/components/lwip/lwip && find ../../../../patches/lwip/ -type f -name '*.patch' -exec patch -p1 -i {} \; )
