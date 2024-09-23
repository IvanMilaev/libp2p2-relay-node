" Load vim-plug plugin manager
call plug#begin('~/.vim/plugged')

" Go development plugin
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

" Auto-completion framework
Plug 'neoclide/coc.nvim', {'branch': 'release'}

call plug#end()

" General Vim settings
syntax on
set number
set tabstop=4
set shiftwidth=4
set expandtab
set autoindent
set smartindent

" Vim-Go settings
let g:go_fmt_command = "goimports"
let g:go_autodetect_gopath = 1

" Set up CoC for Go
let g:coc_global_extensions = ['coc-go']

" Enable code folding
set foldmethod=syntax
set foldlevel=1

" Enable mouse support
set mouse=a

" Map shortcuts for easier navigation and actions
nmap <leader>r :GoRun<CR>
nmap <leader>b :GoBuild<CR>
nmap <leader>t :GoTest<CR>
