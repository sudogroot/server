const gulp = require('gulp');
const sass = require('gulp-sass');
const autoprefixer = require('gulp-autoprefixer');

var paths = {
    styles: {
        src: 'assets/scss/styles.scss',
        all: 'assets/scss/**/**/*.scss',
        dest: 'assets/css/'
    }
}

function styles() {
    return gulp.src(paths.styles.src)
    .pipe(sass().on('error', sass.logError))
    .pipe(autoprefixer({
        browsers: ['last 2 versions'],
        cascade: false
    }))
    .pipe(gulp.dest(paths.styles.dest));
}

function watch() {
    gulp.watch(paths.styles.all, styles);
}

gulp.task('default', gulp.series(watch));