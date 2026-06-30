" ============================================================
" BASICS
" ============================================================
set nocompatible          " forget vi, use vim fully
syntax enable
filetype plugin indent on

set encoding=utf-8
set history=1000

" ============================================================
" APPEARANCE
" ============================================================
set number                " absolute line numbers
set relativenumber        " relative numbers for the rest (so 3j/5k jumps feel natural)
set cursorline            " highlight current line
set colorcolumn=100       " soft ruler at col 100 (Java lines get long)
set scrolloff=8           " keep 8 lines visible above/below cursor
set signcolumn=yes        " always show the gutter (no layout jump when errors appear)

set laststatus=2          " always show status line
set statusline=%f\ %m%r%h%w\ [%Y]\ [%{&ff}]\ %=%l/%L\ col\ %c

set showcmd               " show partial command in bottom right
set showmatch             " briefly jump to matching bracket
set wildmenu              " tab-complete vim commands with a menu
set wildmode=longest:full,full

" ============================================================
" INDENTATION  (Java: 4 spaces, no tabs)
" ============================================================
set tabstop=4
set shiftwidth=4
set softtabstop=4
set expandtab             " spaces not tabs
set smartindent
set autoindent

" ============================================================
" SEARCH
" ============================================================
set incsearch             " search as you type
set hlsearch              " highlight matches
set ignorecase            " case-insensitive by default...
set smartcase             " ...unless you type a capital letter

" clear search highlight with <leader>/ (leader is \ by default)
nnoremap <leader>/ :nohlsearch<CR>

" ============================================================
" FILE HANDLING
" ============================================================
set noswapfile            " no .swp files cluttering the project
set nobackup
set undofile              " persistent undo — survives restarts
set undodir=~/.vim/undo   " store undo files here

set autoread              " reload file if changed outside vim

" ============================================================
" SPLITS  (more natural direction)
" ============================================================
set splitbelow            " horizontal split opens below
set splitright            " vertical split opens to the right

" navigate splits with Ctrl-hjkl
nnoremap <C-h> <C-w>h
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-l> <C-w>l

" ============================================================
" CONVENIENCE MAPS
" ============================================================
" jk to escape insert mode (saves reaching for Esc constantly)
inoremap jk <Esc>

" keep visual selection after indent
vnoremap < <gv
vnoremap > >gv

" move lines up/down in visual mode
vnoremap J :m '>+1<CR>gv=gv
vnoremap K :m '<-2<CR>gv=gv

" Y behaves like D and C (yank to end of line, not whole line)
nnoremap Y y$

" center screen after big jumps
nnoremap n nzzzv
nnoremap N Nzzzv
nnoremap <C-d> <C-d>zz
nnoremap <C-u> <C-u>zz

" open netrw file explorer (built-in, no plugin needed)
nnoremap <leader>e :Lexplore<CR>

" toggle cheatsheet in a right split
nnoremap <leader>? :vsplit ~/vim-cheatsheet.md<CR>

" open terminal in a horizontal split below
nnoremap <leader>t :terminal<CR>

" quick compile + run for a single Java file (useful early chapters)
nnoremap <leader>jc :!javac %<CR>
nnoremap <leader>jr :!java %:r<CR>

" ============================================================
" NETRW  (built-in file browser)
" ============================================================
let g:netrw_banner    = 0   " hide the banner
let g:netrw_liststyle = 3   " tree view
let g:netrw_winsize   = 25  " 25% width sidebar

" ============================================================
" JAVA SPECIFICS
" ============================================================
" tell vim where javadoc ends (helps with /* */ block folding)
autocmd FileType java setlocal foldmethod=syntax foldlevel=99
