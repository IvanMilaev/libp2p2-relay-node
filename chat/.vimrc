" Enable syntax highlighting and line numbers
syntax on
set number

" Set encoding
set encoding=utf-8

" Set tabs and indentation for Python
set tabstop=4           " Number of spaces that a <Tab> counts for
set shiftwidth=4        " Number of spaces to use for auto-indentation
set expandtab           " Use spaces instead of tabs
set autoindent          " Copy indent from the current line when starting a new line
set smartindent         " Automatically insert one extra level of indentation in some cases

" Show matching brackets
set showmatch

" Enable line wrapping and set color column for PEP8 (max line length 79)
set textwidth=79
highlight OverLength ctermbg=darkred ctermfg=white guibg=#592929
match OverLength /\%80v./

" Enable mouse support
set mouse=a

" Enable folding based on indentation
set foldmethod=indent
set foldlevel=99        " Start unfolded

" Automatically save undo history in a file
set undofile
set undodir=~/.vim/undodir

" Enable persistent clipboard between Vim and system
set clipboard=unnamedplus

" Python-specific settings
autocmd FileType python setlocal expandtab shiftwidth=4 tabstop=4 softtabstop=4

" Use the following plugin manager (vim-plug)
call plug#begin('~/.vim/plugged')

" Python syntax checking and linting (install flake8)
Plug 'dense-analysis/ale'  " Asynchronous Lint Engine for linting and fixing
let g:ale_linters = {'python': ['flake8']}
let g:ale_fixers = {'python': ['autopep8', 'black']}
let g:ale_fix_on_save = 1   " Automatically fix code on save

" Python auto-completion
Plug 'davidhalter/jedi-vim'  " Python autocompletion using Jedi

" Git integration
Plug 'tpope/vim-fugitive'  " Git wrapper

" Statusline plugin for better visibility
Plug 'vim-airline/vim-airline'
Plug 'vim-airline/vim-airline-themes'

" Fuzzy finder for files, buffers, etc. (requires fzf)
Plug 'junegunn/fzf'
Plug 'junegunn/fzf.vim'

" Indentation guide
Plug 'Yggdroot/indentLine'

call plug#end()

" Key mappings for convenience
noremap <F8> :ALEFix<CR>         " Press F8 to auto-fix with ALE
noremap <leader>r :w<CR>:!python3 %<CR>  " Run the current Python file with <leader>r (leader is \ by default)

