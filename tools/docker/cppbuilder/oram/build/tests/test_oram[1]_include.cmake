if(EXISTS "/builder/build/tests/test_oram[1]_tests.cmake")
  include("/builder/build/tests/test_oram[1]_tests.cmake")
else()
  add_test(test_oram_NOT_BUILT test_oram_NOT_BUILT)
endif()
