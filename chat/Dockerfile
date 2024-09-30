# Use the official Python image
FROM python:latest

# Install necessary packages
RUN apt-get update && apt-get install -y \
    vim-nox \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the requirements file and install dependencies
COPY requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt

# Install Vim-Plug
RUN curl -fLo /root/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# Copy the .vimrc configuration file
COPY .vimrc /root/.vimrc

# Copy the project files into the container
COPY . /app

# Start an interactive bash shell
CMD ["bash"]

