$(document).ready(function() {
    $("#timeslot-list").hide();
    $("#timeslot-latest").show();
    var seasonTable = $("#seasons").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "pageLength": 5,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
    });
    var timeslotTable = $("#timeslots").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
        "pageLength": 5,
        "columnDefs": [
        {
            "targets": [ 0 ],
            "visible": false,
            "searchable": true
        },
        {
            "targets": [ 1,2 ],
            "searchable": false
        },

        ]
    });

    $.urlParam = function (name) {
        var results = new RegExp('[\?&]' + name + '=([^&#]*)').exec(window.location.href);
        if (results == null) {
            return null;
        }
        else {
            return decodeURI(results[1]) || 0;
        }
    }

    var filterSeason = function(object) {
        $(".timeslot-filter-season").removeClass("active");
        $(object).addClass("active");
        $("#timeslot-latest").hide();
        $("#timeslot-season-title").text("Season " + $(object).find("td:nth-of-type(2)").html())
        timeslotTable
            .columns( 0 )
            .search( $(object).data("seasonid") )
            .draw();

        $("#timeslot-list").show();

        $("#timeslot-list tr").on("click", function () {
            window.location.href = $(object).find("a").attr("href");
        });
    };

    $(".timeslot-filter-season").on( "click", function () {
        filterSeason(this);
    });

    var getSeasonFilter = function() {
        if ($.urlParam("seasonID") != null) {
            var rowId = $("#seasons").dataTable()
                .fnFindCellRowIndexes($.urlParam("seasonID"), 0);
            var seasonTable = $("#seasons").DataTable()
                .cell(7, 1)
                .data('LOL')
                .row(rowId).scrollTo();
            //filterSeason($("tr[data-seasonid='" + $.urlParam("seasonID") + "']"));
            //window.location.href = window.location.href.split('#')[0] + "#timeslot-list";
        }
    }
    getSeasonFilter();
});
