# Use the official Python image as the base image
FROM python:latest

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Set the working directory inside the container
WORKDIR /app

# Install vim-nox, git, curl, and other necessary tools
RUN apt-get update && apt-get install -y \
    vim-nox \
    git \
    curl \
    openssh-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Vim-Plug (Vim plugin manager)
RUN curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# Create a .vimrc file with Vim-Plug configuration for Python development
RUN echo "\
\" Enable syntax highlighting\n\
syntax on\n\
\n\
\" Use 4 spaces for indentation\n\
set tabstop=4\n\
set shiftwidth=4\n\
set expandtab\n\
\n\
\" Enable line numbers\n\
set number\n\
\n\
\" Enable auto-indentation\n\
set autoindent\n\
set smartindent\n\
\n\
\" Enable line wrapping\n\
set wrap\n\
\n\
\" Show matching brackets\n\
set showmatch\n\
\n\
\" Enable mouse support\n\
set mouse=a\n\
\n\
\" Enable relative line numbers\n\
set relativenumber\n\
\n\
\" Use system clipboard\n\
set clipboard=unnamedplus\n\
\n\
\" Set Python-specific settings\n\
autocmd FileType python setlocal expandtab shiftwidth=4 softtabstop=4\n\
autocmd FileType python setlocal colorcolumn=79\n\
\n\
\" Enable file type plugins\n\
filetype plugin on\n\
\n\
\" Use Vim-Plug to manage plugins\n\
call plug#begin('~/.vim/plugged')\n\
\n\
\" Python autocompletion\n\
Plug 'davidhalter/jedi-vim'\n\
\n\
\" File explorer\n\
Plug 'preservim/nerdtree'\n\
\n\
\" Git integration\n\
Plug 'airblade/vim-gitgutter'\n\
\n\
\" Status/tab line\n\
Plug 'vim-airline/vim-airline'\n\
\n\
\" End of plugins\n\
call plug#end()\n\
\n\
\" NERDTree settings\n\
map <C-n> :NERDTreeToggle<CR>\n\
\n\
\" Jedi settings\n\
let g:jedi#completions_enabled = 1\n\
let g:jedi#use_tabs_not_buffers = 1\n\
\n\
\" GitGutter settings\n\
let g:gitgutter_enabled = 1\n\
\n\
\" Airline settings\n\
let g:airline#extensions#tabline#enabled = 1\n\
" > ~/.vimrc

# Copy SSH keys into the container (keys should be available in the context)
# The 'id_rsa' and 'id_rsa.pub' files should be in the same directory as the Dockerfile
COPY git_milaevivan /root/.ssh/id_rsa
COPY git_milaevivan.pub /root/.ssh/id_rsa.pub

# Set appropriate permissions for SSH keys
RUN chmod 600 /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa.pub

# Add GitHub to known hosts to avoid SSH prompt on the first connection
RUN touch /root/.ssh/known_hosts && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the project files into the container
COPY . .

# Expose the port for the relay node (optional)
EXPOSE 10000

# Command to start the relay node
CMD ["python", "relay_node.py"]

