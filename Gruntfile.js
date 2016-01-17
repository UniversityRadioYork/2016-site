module.exports = function (grunt) {

	// Project configuration.
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		wiredep: {
			task: {
				src: [
					'views/**/*.mustache',   // .html support...
				],
				overrides: {
					"bootstrap": {
						"main": [
							"less/bootstrap.less",
							"dist/css/bootstrap.css",
							"dist/js/bootstrap.js"
						]
					}
				},
				ignorePath: '../../public'
			}
		},
		sass: {
			options: {
				sourceMap: true,
				outputStyle: 'compressed'
			},
			dist: {
				files: {
					'public/styles/main.css': 'src/styles/main.scss'
				}
			}
		},
		uglify: {
			main: {
				files: [
					{
						expand: true,
						cwd: 'src/scripts/',
						src: '**/*.js',
						dest: 'public/scripts/',
						ext: '.min.js'
					}
				],
				options: {
					sourceMap: true
				}
			}
		},
		watch: {
			stylesheets: {
				files: ['src/**/*.scss'],
				tasks: ['sass'],
				options: {
					livereload: true
				}
			},
			scripts: {
				files: ['src/**/*.js'],
				tasks: ['uglify'],
				options: {
					livereload: true
				}
			},
			pages: {
				files: ['views/**/*.mustache'],
				tasks: ['wiredep'],
				options: {
					livereload: true
				}
			},
		},
		clean: {
			main: {
				src: [
					"public/scripts/",
					"public/styles/"
				]
			}
		},
		auto_install: {
			main: {
				local: {}
			}
		}
	});

	grunt.loadNpmTasks('grunt-wiredep');
	grunt.loadNpmTasks('grunt-sass');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-auto-install');

	// Default task(s).
	grunt.registerTask('default', ['build', 'watch']);

	// Just for compiling things
	grunt.registerTask('build', ['clean', 'build:noclean']);

	grunt.registerTask('build:noclean', ['wiredep', 'sass', 'uglify']);

};
